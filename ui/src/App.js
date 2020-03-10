
import React, { useState, useEffect } from 'react';
import axios from "axios";
import { Router } from 'react-router-dom';
import { createBrowserHistory } from 'history';
import { ThemeProvider } from '@material-ui/styles';

import theme from './theme';
import './assets/scss/index.scss';
import './App.css';
import Routes from './Routes';
import configureAppStore from './store'
import { Provider } from 'react-redux'
const browserHistory = createBrowserHistory();

const store = configureAppStore({
  inputs: {
    installGit: false,
    installKubeseal: false,
    githubUsernameRegisterd: false,
    cloneTelarWeb: false,
    cloneTsServerless: false,
    cloneTsUi: false,
    openFaaSApp: false,
    openFaaSAppHasRepos: false,
    githubSSHKey: false,
    firebaseStorage: false,
    mongoDBConnection: false,
    websocketConnection: false,
    deployTelarWeb: false,
    deployTsServerless: false,
    deploySocialUi: false
  }, setupState: 'init', setupStep: 0
})

function App() {
  // Define storage for data
  const [state, setState] = useState([]);

  // Retrieving data from an AJAX request.
  // Remember that the function passed to useEffect will run,
  // after render is fixed on the screen.
  // See https://reactjs.org/docs/hooks-reference.html#useeffect
  useEffect(() => {
    axios
      .get("/hello") // GET request to URL /hello
      .then(resp => setState(resp.data)) // save response to state
      .catch(err => console.log(err)); // catch error
  });

  return (
    <Router history={browserHistory}>
      <Provider store={store}>
        <ThemeProvider theme={theme}>
          <Routes />
        </ThemeProvider>
      </Provider>
    </Router>
  );
}

export default App;
