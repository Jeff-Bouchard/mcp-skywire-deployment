import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {DataSource} from '@angular/cdk/collections';
import {BehaviorSubject, Observable} from 'rxjs/Rx';
import {MatPaginator, MatSort} from '@angular/material';
import {ProxyDataService} from '../../services/proxy.data.service';
import {HttpClient} from '@angular/common/http';
import {HttpService} from '../../services/http.service';
import {ProxyService} from '../../models/proxy';
import {animate, state, style, transition, trigger} from '@angular/animations';

@Component({
  selector: 'app-all-configurations',
  templateUrl: './all-configurations.component.html',
  styleUrls: ['./all-configurations.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({height: '0px', minHeight: '0', display: 'none'})),
      state('expanded', style({height: '*'})),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class AllConfigurationsComponent implements OnInit {
  public copyText = 'PROXY.COPY';

  dataSource: ExampleDataSource | null;
  columnsToDisplay = ['keyPort', 'country', 'city', 'status'];

  constructor(public httpClient: HttpClient, public httpService: HttpService, public transactionDataService: ProxyDataService) { }

  @ViewChild(MatPaginator) paginator: MatPaginator;
  @ViewChild(MatSort) sort: MatSort;
  @ViewChild('filter') filter: ElementRef;

  ngOnInit() {
    this.loadData(this);
  }

  public copyToClipboard(secretInfo: string): void {
    let body = document.getElementsByTagName('body')[0];
    const input: any = document.createElement('INPUT');
    body.appendChild(input);
    input.setAttribute('value', secretInfo);
    const isiOSDevice = navigator.userAgent.match(/ipad|iphone/i);

    if (isiOSDevice) {
      var editable = input.contentEditable;
      var readOnly = input.readOnly;

      input.contentEditable = true;
      input.readOnly = false;

      var range = document.createRange();
      range.selectNodeContents(input);

      var selection = window.getSelection();
      selection.removeAllRanges();
      selection.addRange(range);

      input.setSelectionRange(0, 999999);
      input.contentEditable = editable;
      input.readOnly = readOnly;
    } else {
      input.select();
    }
    document.execCommand('copy');
    body.removeChild(input);
    this.copyText = 'PROXY.COPIED';
  }

  public changeText(): void {
    this.copyText = 'PROXY.COPY';
  }

  public loadData(_this) {
    _this.dataService = new ProxyDataService(_this.httpClient);
    _this.dataSource = new ExampleDataSource(_this.dataService, _this.paginator, _this.sort);
    Observable.fromEvent(_this.filter.nativeElement, 'keyup')
      .debounceTime(150)
      .distinctUntilChanged()
      .subscribe(() => {
        if (!_this.dataSource) {
          return;
        }
        let stringValue = _this.filter.nativeElement.value.toLowerCase();
        _this.dataSource.filter = stringValue;
      });
  }
}


export class ExampleDataSource extends DataSource<ProxyService> {
  _filterChange = new BehaviorSubject('');
  _filterStatus = new BehaviorSubject('');

  get filter(): string {
    return this._filterChange.value;
  }

  set filter(filter: string) {
    this._filterChange.next(filter);
  }

  get activeState(): string {
    return this._filterStatus.value;
  }

  set activeState(activeState: string) {
    this._filterStatus.next(activeState);
  }

  filteredData: any[] = [];
  renderedData: any[] = [];

  constructor(public _exampleDatabase: ProxyDataService,
              public _paginator: MatPaginator,
              public _sort: MatSort) {
    super();
    // Reset to the first page when the user changes the filter.
    this._filterChange.subscribe(() => this._paginator.pageIndex = 0);
  }

  /** Connect function called by the table to retrieve one stream containing the data to render. */
  connect(): Observable<ProxyService[]> {
    // Listen for any changes in the base data, sorting, filtering, or pagination
    const displayDataChanges = [
      this._exampleDatabase.dataChange,
      this._sort.sortChange,
      this._filterChange,
      this._filterStatus,
      this._paginator.page
    ];

    this._exampleDatabase.getAllProxies();

    return Observable.merge(...displayDataChanges).map(() => {
      // Filter data
      this.filteredData = this._exampleDatabase.data.slice().filter((proxy: ProxyService) => {
        const searchStr = (proxy.publicKeyPort).toLowerCase();
        return searchStr.indexOf(this.filter.toLowerCase()) !== -1;
      });

      // Sort filtered data
      const sortedData = this.sortData(this.filteredData.slice());

      // Grab the page's slice of the filtered sorted data.
      const startIndex = this._paginator.pageIndex * this._paginator.pageSize;
      this.renderedData = sortedData.splice(startIndex, this._paginator.pageSize);
      return this.renderedData;
    });
  }
  disconnect() {
  }



  /** Returns a sorted copy of the database data. */
  sortData(data: ProxyService[]): ProxyService[] {
    if (!this._sort.active || this._sort.direction === '') {
      return data;
    }

    return data.sort((a, b) => {
      let propertyA: number | string = '';
      let propertyB: number | string = '';

      switch (this._sort.active) {
        case 'keyPort': [propertyA, propertyB] = [a.publicKeyPort, b.publicKeyPort]; break;
        case 'country': [propertyA, propertyB] = [a.country, b.country]; break;
        case 'city': [propertyA, propertyB] = [a.city, b.city]; break;
        case 'status': [propertyA, propertyB] = [a.status, b.status]; break;
      }

      const valueA = isNaN(+propertyA) ? propertyA : +propertyA;
      const valueB = isNaN(+propertyB) ? propertyB : +propertyB;

      return (valueA < valueB ? -1 : 1) * (this._sort.direction === 'asc' ? 1 : -1);
    });
  }
}

