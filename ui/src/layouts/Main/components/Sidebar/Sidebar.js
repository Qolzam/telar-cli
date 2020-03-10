import React from 'react';
import clsx from 'clsx';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/styles';
import { Drawer } from '@material-ui/core';
import SetupIcon from '@material-ui/icons/Build';
import { withRouter, matchPath } from 'react-router-dom';
import { SidebarNav } from './components';
import SettingsEthernetIcon from '@material-ui/icons/SettingsEthernet';
import DeviceIcon from '@material-ui/icons/DevicesOther';
import ChartIcon from '@material-ui/icons/InsertChartOutlined';

const useStyles = makeStyles(theme => ({
  drawer: {
    width: 240,
    [theme.breakpoints.up('lg')]: {
      marginTop: 64,
      height: 'calc(100% - 64px)'
    }
  },
  root: {
    backgroundColor: theme.palette.white,
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
    padding: theme.spacing(2)
  },
  divider: {
    margin: theme.spacing(2, 0)
  },
  nav: {
    marginBottom: theme.spacing(2)
  }
}));

const Sidebar = props => {
  const { open, variant, onClose, className, ...rest } = props;
  let user = ''
  const { pathname } = props.history.location;
  const matchFunction = matchPath(pathname, {
    path: '/:user',
    strict: true,
  });

  if(matchFunction) {
    user = matchFunction.params.user
  }
  const classes = useStyles();

  const pages = [
    {
      title: 'Setup Telar Social',
      href: `/setup`,
      icon: <SetupIcon />
    }
  ];

  return (
    <Drawer
      anchor="left"
      classes={{ paper: classes.drawer }}
      onClose={onClose}
      open={open}
      variant={variant}
    >
      <div
        {...rest}
        className={clsx(classes.root, className)}
      >
        <SidebarNav
          className={classes.nav}
          pages={pages}
        />
      </div>
    </Drawer>
  );
};

Sidebar.propTypes = {
  className: PropTypes.string,
  onClose: PropTypes.func,
  open: PropTypes.bool.isRequired,
  variant: PropTypes.string.isRequired
};

export default withRouter(Sidebar);
