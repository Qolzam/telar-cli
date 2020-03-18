import React from 'react';
import { useDispatch, useSelector } from 'react-redux'
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormGroup from '@material-ui/core/FormGroup';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';
import CircularProgress from '@material-ui/core/CircularProgress';
import { makeStyles, withStyles } from '@material-ui/core/styles';
import { green } from '@material-ui/core/colors';

const GreenCheckbox = withStyles({
  root: {
    color: green[400],
    '&$checked': {
      color: green[600],
    },
  },
  checked: {},
})(props => <Checkbox color="default" {...props} />);

const useStyles = makeStyles({
    root: {
      minWidth: 275,
    },
    bullet: {
      display: 'inline-block',
      margin: '0 2px',
      transform: 'scale(0.8)',
    },
    title: {
      fontSize: 14,
    },
    pos: {
      marginBottom: 12,
    },
    progress: {
      margin: '10px'
    }
  });
  
export default function DeployDialog(props) {
  const classes = useStyles();
  const {open, onClose} = props

  const loadingStackYaml = useSelector(state => state['inputs']['loadingStackYaml'])
  const loadingCreateSecret = useSelector(state => state['inputs']['loadingCreateSecret'])
  const loadingPublicPrivateKey = useSelector(state => state['inputs']['loadingPublicPrivateKey'])
  const deployTelarWeb = useSelector(state => state['inputs']['deployTelarWeb'])
  const deployTsServerless = useSelector(state => state['inputs']['deployTsServerless'])
  const deploySocialUi = useSelector(state => state['inputs']['deploySocialUi'])

  const handleClose = () => {
    onClose()
  };
  const checkBox = (checked) => {
    if (checked) {
      return <GreenCheckbox checked={true} />
    }
    return <CircularProgress className={classes.progress} size={17} color="secondary" />
  }

  return (
      <Dialog
        open={open}
        disableBackdropClick={true}
        disableEscapeKeyDown={true}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">{"Deploying"}</DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            Deploying the project to OpenFaaS Cloud Community Cluster
          </DialogContentText>
            <FormControl component="fieldset" className={classes.formControl}>
            <FormGroup>


            <FormControlLabel
              control={checkBox(loadingStackYaml)}
              checked={loadingStackYaml}
              label="Creating Configurations"
            />

            <FormControlLabel
              control={checkBox(loadingCreateSecret)}
              checked={loadingCreateSecret}
              label="Creating Secrets"
            />

            <FormControlLabel
              control={checkBox(loadingPublicPrivateKey)}
              checked={loadingPublicPrivateKey}
              label="Deploying Telar Web"
            />

            <FormControlLabel
              control={checkBox(deployTelarWeb)}
              checked={deployTelarWeb}
              label="Deploying Telar Web"
            />

            <FormControlLabel
              control={checkBox(deployTsServerless)}
              checked={deployTsServerless}
              label="Deploying Telar Social Serverless"
            />

            <FormControlLabel
              control={checkBox(deploySocialUi)}
              checked={deploySocialUi}
              label="Deploying Telar Social User Interface"
            />
           
           
            </FormGroup>
            </FormControl>
        </DialogContent>
        <DialogActions>
        </DialogActions>
      </Dialog>
  );
}