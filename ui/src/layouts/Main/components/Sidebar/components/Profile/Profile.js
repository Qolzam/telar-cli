import React from 'react';
import { Link as RouterLink } from 'react-router-dom';
import clsx from 'clsx';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/styles';
import { Avatar, Typography } from '@material-ui/core';
import { withRouter, matchPath } from 'react-router-dom';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    minHeight: 'fit-content'
  },
  avatar: {
    width: 60,
    height: 60
  },
  name: {
    marginTop: theme.spacing(1)
  }
}));

const Profile = props => {
  const { className, ...rest } = props;
  let userLogin = 'de-amir'
  const { pathname } = props.history.location;
  const matchFunction = matchPath(pathname, {
    path: '/:user',
    strict: true,
  });

  if(matchFunction) {
    userLogin =  matchFunction.params.user.toLowerCase()
  }
  const classes = useStyles();
  const users = {
    'de-amir': { 
      name: 'Amirhossein Movahedi',
      avatar: '/images/avatars/amir.JPG',
      bio: 'Distributed System Eng.'
    },
    'de-minhnguyen': { 
      name: 'Minh Tuan Nguyen',
      avatar: '/images/avatars/minh.JPG',
      bio: 'AI Eng.'
    }
  };
  const user = users['de-amir']

  return (
    <div
      {...rest}
      className={clsx(classes.root, className)}
    >
      <Avatar
        alt="Person"
        className={classes.avatar}
        component={RouterLink}
        src={user.avatar}
        to="/settings"
      />
      <Typography
        className={classes.name}
        variant="h4"
      >
        {user.name}
      </Typography>
      <Typography variant="body2">{user.bio}</Typography>
    </div>
  );
};

Profile.propTypes = {
  className: PropTypes.string
};

export default withRouter(Profile);
