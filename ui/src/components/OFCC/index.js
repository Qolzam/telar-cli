import React from 'react';
import { useDispatch, useSelector } from 'react-redux'
import { makeStyles } from '@material-ui/core/styles';
import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepLabel from '@material-ui/core/StepLabel';
import StepContent from '@material-ui/core/StepContent';
import Button from '@material-ui/core/Button';
import InitialStep from '../InitialStep';
import OFCInfo from '../OFCInfo';
import CheckIngredients from '../CheckIngredients';
import FirebaseStorage from '../FirebaseStorage';
import UserManagement from '../UserManagement';
import WebSocket from '../Websocket';
import Database from '../Database';
import GoogleReCaptcha from '../GoogleReCaptcha';
import OAuth from '../OAuth';
import DeployDialog from '../DepolyDialog';
import actions from '../../store/actions'
import services from '../../services';
const useStyles = makeStyles(theme => ({
  root: {
    width: '100%',
  },
  button: {
    marginTop: theme.spacing(1),
    marginRight: theme.spacing(1),
  },
  actionsContainer: {
    marginBottom: theme.spacing(2),
  },
  resetContainer: {
    padding: theme.spacing(3),
  },
}));

function getSteps() {
  return [
  `Let's start`, 
  `General settings`, 
  'Check ingredients', 
  'Firebase Storage', 
  'MongoDB', 
  'Google reCAPTCHA', 
  'Github OAuth',
  'User Manegement',
  'Websocket'  
];
}

function getStepContent(step) {
  switch (step) {
    case 0:
      return <InitialStep />;
    case 1:
      return <OFCInfo />;
    case 2:
      return  <CheckIngredients />;
    case 3:
      return  <FirebaseStorage />;
    case 4:
      return  <Database />;
    case 5:
      return  <GoogleReCaptcha />;
    case 6:
      return  <OAuth />;
    case 7:
      return  <UserManagement />;
    case 8:
      return  <WebSocket />;
    default:
      return 'Unknown step';
  }
}

const validInput = (value) => {
  if(value && value.trim() !== "") {
    return true
  }
  return false
}

const validInputs = (state, inputs) => {
  for (const input of inputs) {
    const isValid = validInput(state['inputs'][input])
    if (!isValid) {
      return false
    }
   }
   return true

}

const validCheckbox = (state, inputs) => {
  for (const input of inputs) {
    if (!state['inputs'][input]) {
      return false
    }
  }
   
    return true
}



export default function OFCC() {
  const classes = useStyles();
  const dispatch = useDispatch()
  const setupStep = useSelector(state => state['setupStep'])
  const deployOpen = useSelector(state => state['deployOpen'])
  const stepCondition = {}
  const steps = getSteps();

  
  // ***** Conditions ***** //

 const state = useSelector(state =>  state) 

  // Init step
  stepCondition[0]= (validInputs(state, ['appID', 'appName', 'companyName', 'supportEmail', 'projectDirectory' ]))

  // Init step
  stepCondition[1]= (validInputs(state, ['ofGateway', 'dockerUser', 'ofUsername', 'socialDomain', 'secretName', 'namespace']))
  
  // Check ingredients
  stepCondition[2] = !validCheckbox(state,['loadingCheckIngredients'])

  // Firebase storage
  stepCondition[3] = (!validCheckbox(state, ['loadingFirebaseStorage'] ) && validInputs(state, ['bucketName']) === true)

  // Database
  stepCondition[4] = (!validCheckbox(state, ['loadingMongoDB'] ) && validInputs(state, ['mongoDBURI', 'mongoDBName']) === true)

  // Firebase storage
  stepCondition[5] = (validInputs(state, ['siteKeyRecaptcha', 'recaptchaKey']))

  // OAuth
  stepCondition[6] = (validInputs(state, ['githubOAuthSecret']))

  // User management
  stepCondition[7] = (validInputs(state, ['adminUsername', 'adminPassword', 'gmail', 'gmailPassword']))

  // Websocket
  stepCondition[8] = ((!validCheckbox(state, ['loadingWebsocket'] ) || validCheckbox(state, ['websocketConnection'] )) && validInputs(state, ['gateway', 'payloadSecret', 'websocketURL']) === true)

  const handleCloseDeploy = () => {
    dispatch(actions.setInput('deployOpen', false))
  };

  const handleNext = () => {

    services.dispatchServer(actions.checkStep(state))
    
  };



  return (
    <div className={classes.root}>
      <Stepper activeStep={setupStep} orientation="vertical" >
        {steps.map((label, index) => (
          <Step key={label}>
            <StepLabel>{label}</StepLabel>
            <StepContent>
              {getStepContent(index)}
              <div className={classes.actionsContainer}>
                <div>
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={handleNext}
                    className={classes.button}
                    disabled={stepCondition[setupStep] !== true}
                  >
                    {setupStep === steps.length - 1 ? 'Deploy' : 'Next'}
                  </Button>
                </div>
              </div>
            </StepContent>
          </Step>
        ))}
      </Stepper>
      

      <DeployDialog open={deployOpen} onClose={handleCloseDeploy} />
    </div>
  );
}