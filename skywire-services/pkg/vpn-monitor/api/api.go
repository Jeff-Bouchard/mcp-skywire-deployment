// Package api pkg/vpn-monitor/api/api.go
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/skycoin/skywire-utilities/pkg/buildinfo"
	"github.com/skycoin/skywire-utilities/pkg/cipher"
	"github.com/skycoin/skywire-utilities/pkg/httputil"
	"github.com/skycoin/skywire-utilities/pkg/logging"
	"github.com/skycoin/skywire/pkg/app/appserver"
	"github.com/skycoin/skywire/pkg/restart"
	"github.com/skycoin/skywire/pkg/servicedisc"
	"github.com/skycoin/skywire/pkg/skyenv"
	"github.com/skycoin/skywire/pkg/transport/network"
	"github.com/skycoin/skywire/pkg/visor"
	"github.com/skycoin/skywire/pkg/visor/visorconfig"

	"github.com/skycoin/skywire-services/internal/vpn"
)

// API register all the API endpoints.
// It implements a net/http.Handler.
type API struct {
	http.Handler
	Config
	ServicesURLs

	Visor *visor.Visor

	vpnKeys   []cipher.PubKey
	deadVPNs  []string
	logger    logging.Logger
	startedAt time.Time
}

// Config is struct for keys and sign value of VM
type Config struct {
	PK   cipher.PubKey
	SK   cipher.SecKey
	Sign cipher.Sig
}

// ServicesURLs is struct for organizing URL's of services
type ServicesURLs struct {
	SD string
	UT string
}

// HealthCheckResponse is struct of /health endpoint
type HealthCheckResponse struct {
	BuildInfo *buildinfo.Info `json:"build_info,omitempty"`
	StartedAt time.Time       `json:"started_at,omitempty"`
}

// Error is the object returned to the client when there's an error.
type Error struct {
	Error string `json:"error"`
}

// New returns a new *chi.Mux object, which can be started as a server
func New(logger *logging.Logger, srvURLs ServicesURLs, vmConfig Config) *API {

	api := &API{
		Config:       vmConfig,
		ServicesURLs: srvURLs,
		logger:       *logger,
		startedAt:    time.Now(),
	}
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httputil.SetLoggerMiddleware(logger))
	r.Get("/health", api.health)
	api.Handler = r

	return api
}

func (api *API) health(w http.ResponseWriter, r *http.Request) {
	info := buildinfo.Get()
	api.writeJSON(w, r, http.StatusOK, HealthCheckResponse{
		BuildInfo: info,
		StartedAt: api.startedAt,
	})
}

func (api *API) writeJSON(w http.ResponseWriter, r *http.Request, code int, object interface{}) {
	jsonObject, err := json.Marshal(object)
	if err != nil {
		api.log(r).WithError(err).Errorf("failed to encode json response")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(jsonObject)
	if err != nil {
		api.log(r).WithError(err).Errorf("failed to write json response")
	}
}

func (api *API) log(r *http.Request) logrus.FieldLogger {
	return httputil.GetLogger(r)
}

// InitDeregistrationLoop is function which runs periodic background tasks of API.
func (api *API) InitDeregistrationLoop(ctx context.Context, conf *visorconfig.V1, sleepDeregistration time.Duration) {
	// Start a visor
	api.startVisor(ctx, conf)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			api.deregister()
			time.Sleep(sleepDeregistration * time.Minute)
		}
	}
}

// deregister dead VPNs entries in service discovery
func (api *API) deregister() {
	api.logger.Info("VPN Deregistration started.")

	// reload keys
	api.getVPNKeys()

	// monitoring VPNs
	onlineVpnCount := int64(0)
	api.deadVPNs = []string{}

	if len(api.vpnKeys) == 0 {
		api.logger.Warn("No VPN keys found")
	} else {
		for _, key := range api.vpnKeys {
			api.testVPN(key, &onlineVpnCount)
		}
		api.logger.WithField("count", onlineVpnCount).Info("VPNs online.")

		// deregister dead VPNs
		if len(api.deadVPNs) > 0 {
			api.vpnDeregister(api.deadVPNs)
		}
	}

	api.logger.WithField("Number of dead VPNs", len(api.deadVPNs)).WithField("PKs", api.deadVPNs).Info("VPN Deregistration completed.")
}

func (api *API) testVPN(key cipher.PubKey, onlineVpnCount *int64) {

	online := api.isOnline(key)

	if online {
		atomic.AddInt64(onlineVpnCount, 1)
	}

	if !online {
		api.deadVPNs = append(api.deadVPNs, key.Hex())
	}
}

func (api *API) isOnline(key cipher.PubKey) (ok bool) {
	transport := network.DMSG

	tp, err := api.Visor.AddTransport(key, string(transport), time.Second*10)
	if err != nil {
		api.logger.WithError(err).Warnf("Failed to establish %v transport", transport)
		return false
	}

	var latency time.Duration
	api.logger.Infof("Established %v transport to %v", transport, key)
	// We use the name vpn-client and not vpn-lite-client here to get around the constraint that
	// -srv flag can only be set for vpn-client and skysocks-client.
	// And due to this the binary should also be named as vpn-client and not vpn-client-lite
	sum, vpnErr := RunVpnClient(api.Visor, key, skyenv.VPNClientName)
	time.Sleep(time.Second * 4)
	ok = true

	switch vpnErr {
	case nil:
		if len(sum) > 0 {
			latency = sum[0].Latency
		}
	case vpn.ErrSetupNode, vpn.ErrNotPermitted:
		api.logger.WithError(vpnErr).Infof("Vpn error on %v transport of %v.", transport, key)
	default:
		api.logger.WithError(vpnErr).Infof("Vpn error on %v transport of %v.", transport, key)
		ok = false
	}

	err = api.Visor.RemoveTransport(tp.ID)
	if err != nil {
		api.logger.Warnf("Error removing %v transport of %v: %v", transport, key, err)
	}

	if ok && latency != 0 {
		return ok
	}

	return ok
}

func (api *API) vpnDeregister(keys []string) {
	err := api.deregisterRequest(keys, fmt.Sprintf(api.ServicesURLs.SD+"/api/services/deregister/vpn"))
	if err != nil {
		api.logger.Warn(err)
		return
	}
	api.logger.Info("Deregister request send to SD")
}

// deregisterRequest is deregistration handler for all services
func (api *API) deregisterRequest(keys []string, rawReqURL string) error {
	reqURL, err := url.Parse(rawReqURL)
	if err != nil {
		return fmt.Errorf("Error on parsing deregistration URL : %v", err)
	}

	jsonData, err := json.Marshal(keys)
	if err != nil {
		return fmt.Errorf("Error on parsing deregistration keys : %v", err)
	}
	body := bytes.NewReader(jsonData)

	req := &http.Request{
		Method: "DELETE",
		URL:    reqURL,
		Header: map[string][]string{
			"NM-PK":   {api.Config.PK.Hex()},
			"NM-Sign": {api.Config.Sign.Hex()},
		},
		Body: io.NopCloser(body),
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error on send deregistration request : %s", err)
	}
	defer res.Body.Close() //nolint

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deregistering vpn keys: status code %v", res.StatusCode)
	}

	return nil
}

type vpnList []servicedisc.Service

func getVPNs(sdURL string) (data vpnList, err error) {
	res, err := http.Get(sdURL + "/api/services?type=vpn") //nolint

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (api *API) getVPNKeys() {
	vpns, err := getVPNs(api.ServicesURLs.SD)
	if err != nil {
		api.logger.Warn("Error while fetching vpns: %v", err)
		return
	}
	if len(vpns) == 0 {
		api.logger.Warn("No vpns found... Trying again")
	}
	api.vpnKeys = []cipher.PubKey{}
	for _, vpn := range vpns {
		api.vpnKeys = append(api.vpnKeys, vpn.Addr.PubKey())
	}

	api.logger.WithField("vpns", len(vpns)).Info("Vpn keys updated.")
}

func (api *API) startVisor(ctx context.Context, conf *visorconfig.V1) {
	conf.SetLogger(logging.NewMasterLogger())
	v, ok := visor.NewVisor(ctx, conf, restart.CaptureContext(), false, "", "")
	if !ok {
		api.logger.Fatal("Failed to start visor.")
	}
	api.Visor = v
}

// RunVpnClient runs a vpn-client which connects to vpn server
func RunVpnClient(v *visor.Visor, serverPK cipher.PubKey, appName string) ([]appserver.ConnectionSummary, error) {
	err := v.SetAppPK(appName, serverPK)
	if err != nil {
		return []appserver.ConnectionSummary{}, err
	}
	err = v.StartApp(appName)
	if err != nil {
		return []appserver.ConnectionSummary{}, err
	}

	time.Sleep(time.Second * 15)
	appErr, _ := v.GetAppError(appName) //nolint
	if appErr == vpn.ErrSetupNode.Error() {
		return []appserver.ConnectionSummary{}, vpn.ErrSetupNode
	}
	if appErr == vpn.ErrNotPermitted.Error() {
		return []appserver.ConnectionSummary{}, vpn.ErrNotPermitted
	}
	if appErr == vpn.ErrServerOffline.Error() {
		return []appserver.ConnectionSummary{}, vpn.ErrServerOffline
	}
	sum, err := v.GetAppConnectionsSummary(appName)
	if err != nil {
		return []appserver.ConnectionSummary{}, err
	}
	time.Sleep(time.Second * 2)
	err = v.StopApp(appName)
	if err != nil {
		return []appserver.ConnectionSummary{}, err
	}

	return sum, nil
}
