import React from 'react';
import { useDispatch, useSelector } from 'react-redux'
import { makeStyles } from '@material-ui/core/styles';
import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepLabel from '@material-ui/core/StepLabel';
import StepContent from '@material-ui/core/StepContent';
import Button from '@material-ui/core/Button';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import GithubUsername from '../InitialStep';
import CheckIngredients from '../CheckIngredients';
import FirebaseStorage from '../FirebaseStorage';
import UserManagement from '../UserManagement';
import WebSocket from '../Websocket';
import Database from '../Database';
import GoogleReCaptcha from '../GoogleReCaptcha';
import OAuth from '../OAuth';
import DeployDialog from '../DepolyDialog';
import actions from '../../store/actions'
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
      return <GithubUsername />;
    case 1:
      return  <CheckIngredients />;
    case 2:
      return  <FirebaseStorage />;
    case 3:
      return  <Database />;
    case 4:
      return  <GoogleReCaptcha />;
    case 5:
      return  <OAuth />;
    case 6:
      return  <UserManagement />;
    case 7:
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
  const stepCondition = {}
  const steps = getSteps();


  // ***** Conditions ***** //

  // Init step
  stepCondition[0]= useSelector(state => validInputs(state, ['githubUsername','projectDirectory']))
  
  // Check ingredients
  stepCondition[1] =useSelector(state =>  validCheckbox(state, ['installKubeseal', 'githubUsernameRegisterd', 'cloneTelarWeb', 'cloneTsServerless', 'cloneTsUi', 'openFaaSApp', 'openFaaSAppHasRepos', 'githubSSHKey','projectDirectory' ]))

  // Firebase storage
  stepCondition[2] = useSelector(state =>  validCheckbox(state, ['firebaseStorage'] ) === true && validInputs(state, ['bucketName']) === true)

  // Database
  stepCondition[3] = useSelector(state =>  validCheckbox(state, ['mongoDBConnection'] ) === true && validInputs(state, ['mongoDBUsername', 'mongoDBPassword', 'mongoDBName']) === true)

  // Firebase storage
  stepCondition[5] = useSelector(state =>  validInputs(state, ['siteKeyRecaptcha', 'recaptchaKey']))

  // OAuth
  stepCondition[5] = useSelector(state =>  validInputs(state, ['githubOAuthSecret']))

  // User management
  stepCondition[6] = useSelector(state =>  validInputs(state, ['adminUsername', 'adminPassword', 'gmail', 'gmailPassword']))

  // Websocket
  stepCondition[7] = useSelector(state =>  validCheckbox(state, ['websocketConnection'] ) === true && validInputs(state, ['gateway', 'payloadSecret', 'websocketURL']) === true)


  const [deployOpen, setDeployOpen] = React.useState(false);

  const handleDeploy = () => {
    setDeployOpen(true);
  };

  const handleCloseDeploy = () => {
    setDeployOpen(false);
  };

  const handleNext = () => {
    if (setupStep === 7) {
      setDeployOpen(true)
      return
    }
    dispatch(actions.setSetupStep(setupStep + 1))
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