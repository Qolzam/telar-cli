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

export default function OFCC() {
  const classes = useStyles();
  const dispatch = useDispatch()
  const setupStep = useSelector(state => state['setupStep'])
  const steps = getSteps();

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