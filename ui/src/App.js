
import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux'
import { Router } from 'react-router-dom';
import { createBrowserHistory } from 'history';
import { ThemeProvider } from '@material-ui/styles';
import Snackbar from '@material-ui/core/Snackbar';
import IconButton from '@material-ui/core/IconButton';
import CloseIcon from '@material-ui/icons/Close';
import theme from './theme';
import './assets/scss/index.scss';
import './App.css';
import configureAppStore from './store'
import { Provider } from 'react-redux'
import actions from './store/actions';
import DialogInfo from './components/DialogInfo'
import Master from './Master';
import Routes from "./Routes";

const browserHistory = createBrowserHistory();



const store = configureAppStore({
  inputs: {
    baseAPIRoute: '',
    baseHref: '',
    appName: 'Telar Social Engine',
    companyName: 'Telar',
    supportEmail: 'support@telar.dev',
    appID: 'my-social-network',
    dockerUser: '',
    ofUsername: 'admin',
    ofGateway: 'openfaas.example.com',
    socialDomain: 'https://mysocial.com',
    secretName: 'secrets',
    namespace: 'openfaas-fn',
    projectDirectory: "",
    bucketName: "",
    mongoDBURI: "",
    mongoDBName: "test",
    siteKeyRecaptcha: "",
    recaptchaKey: "",
    githubOAuthSecret: "",
    githubOAuthClientID: "",
    adminUsername: "",
    adminPassword: "",
    gmail: "",
    gmailPassword: "",
    gateway: "",
    payloadSecret: "",
    websocketURL: "",
    installGit: false,
    installKubeseal: false,
    cloneTelarWeb: false,
    cloneTsServerless: false,
    cloneTsUi: false,
    openFaaSApp: false,
    openFaaSAppHasRepos: false,
    githubSSHKey: false,
    firebaseServiceAccount: false,
    firebaseStorage: false,
    mongoDBConnection: false,
    websocketConnection: false,
    loadingCheckIngredients: false,
    loadingFirebaseStorage: false,
    loadingMongoDB: false,
    loadingWebsocket: false,
    loadingStackYaml: false,
    loadingCreateSecret: false,
    loadingPublicPrivateKey: false,
    deployTelarWeb: false,
    deployTsServerless: false,
    deploySocialUi: false,
    deployOpen: false,
  }, setupState: 'init', setupStep: 0, stepCondition: {}, infoDialog: {message: "", url: "", open: false}
})

let ws;

if (window.WebSocket === undefined) {
    alert("Your browser does not support WebSockets");
   
} else {
    ws = initWS();
}

function initWS() {
    let socket = new WebSocket("ws://localhost:31115/ws")
     
    socket.onopen = function() {
       store.dispatch(actions.popMessage("Ready!"))
    };
    socket.onmessage = function (e) {
      const action = JSON.parse(e.data)
      console.log('action', action)
      console.log('e.data', e.data)
      store.dispatch(action)
    }
    socket.onclose = function () {
      store.dispatch(actions.popMessage("Socket closed"))
    }

    return socket;
}

function AppSnackbar() {
  const dispatch = useDispatch()

  const popMessage = useSelector(state => state['popMessage'])

  const handleClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }

    dispatch(actions.popMessage(""))
  };

  return (
    <Snackbar
    anchorOrigin={{
      vertical: 'bottom',
      horizontal: 'left',
    }}
    open={popMessage !== ""}
    autoHideDuration={6000}
    onClose={handleClose}
    message={popMessage}
    action={
      <React.Fragment>
        <IconButton size="small" aria-label="close" color="inherit" onClick={handleClose}>
          <CloseIcon fontSize="small" />
        </IconButton>
      </React.Fragment>
    }
  />
  )
}


function App() {

  // Retrieving data from an AJAX request.
  // Remember that the function passed to useEffect will run,
  // after render is fixed on the screen.
  // See https://reactjs.org/docs/hooks-reference.html#useeffect
  useEffect(() => {
  });

  
  return (
    <Router history={browserHistory}>
      <Provider store={store}>
        <ThemeProvider theme={theme}>
          <Master />
          <Routes />
          <AppSnackbar />
          <DialogInfo />
        </ThemeProvider>
      </Provider>
    </Router>
  );
}

export default App;
