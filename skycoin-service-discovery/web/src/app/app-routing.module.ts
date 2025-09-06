import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {AllConfigurationsComponent} from './components/all-configurations/all-configurations.component';
import {LayoutComponent} from './layout/layout.component';
import {PageNotFoundComponent} from './page-not-found/page-not-found.component';

const routes: Routes = [
  {
    path: '',
    component: LayoutComponent,
    children: [
      {path: '', component: AllConfigurationsComponent}
    ]
  },
  {path: '**', component: PageNotFoundComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
