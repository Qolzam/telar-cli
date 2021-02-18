import React from 'react';
import { Switch, Redirect } from 'react-router-dom';

import { RouteWithLayout } from './components';
import { Main as MainLayout, Minimal as MinimalLayout } from './layouts';
import {
  NotFound as NotFoundView
} from './views';

import SetupComponent from './pages/Setup';
import SettingsComponent from './pages/Settings';

const Routes = () => {
  return (
    <Switch>
      <Redirect
        exact
        from="/"
        to="/setup"
      />
      <RouteWithLayout
        component={SetupComponent}
        exact
        layout={MainLayout}
        path="/setup"
      />
      <RouteWithLayout
        component={SettingsComponent}
        exact
        layout={MainLayout}
        path="/settings"
      />
      <RouteWithLayout
        component={NotFoundView}
        exact
        layout={MinimalLayout}
        path="/not-found"
      />
      <Redirect to="/not-found" />
    </Switch>
  );
};

export default Routes;
