import React from 'react';
import { useDispatch, useSelector } from 'react-redux'
import { makeStyles, withStyles } from '@material-ui/core/styles';
import { green } from '@material-ui/core/colors';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import FormLabel from '@material-ui/core/FormLabel';
import FormControl from '@material-ui/core/FormControl';
import FormGroup from '@material-ui/core/FormGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormHelperText from '@material-ui/core/FormHelperText';
import Checkbox from '@material-ui/core/Checkbox';
import CircularProgress from '@material-ui/core/CircularProgress';
import actions from '../../store/actions'
import HelpDialog from '../HelpDialog';
import services from '../../services';

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

export default function Websocket() {
  const classes = useStyles();
  const dispatch = useDispatch()
  const gateway = useSelector(state => state['inputs']['gateway'])
  const payloadSecret = useSelector(state => state['inputs']['payloadSecret'])
  const websocketURL = useSelector(state => state['inputs']['websocketURL'])
  const websocketConnection = useSelector(state => state['inputs']['websocketConnection'])

  const loadingWebsocket = useSelector(state => state['inputs']['loadingWebsocket'])

  const bull = <span className={classes.bullet}>•</span>;
  const [state, setState] = React.useState({
    gilad: true,
    jason: false,
    antoine: false,
  });
  const [helpOpen, setHelpOpen] = React.useState(false);

  const handleHelp = () => {
   services.openURL("https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/8.md")
  };

  const handleCloseHelp = () => {
    setHelpOpen(false);
  };
  
  const helpContent = () => (
    <Card className={classes.root} variant="outlined">
    <CardContent>
      <Typography variant={'h5'} color="textPrimary" gutterBottom>
      In this step you need to enter your Github account username and a directory path where you want to keep Telar Social project source code. 
  </Typography>
  <br/>
    </CardContent>
    <CardActions>
    </CardActions>
  </Card>
  )
  const handleChange = name => event => {
    dispatch(actions.setInput(name,event.currentTarget.value))
  };
  const checkBox = (checked) => {
    if (checked) {
      return <GreenCheckbox checked={true} />
    }
    return <CircularProgress className={classes.progress} size={17} color="secondary" />
  }

  return (
    <Card className={classes.root} variant="outlined">
      <CardContent>
      <Button variant="outlined" color="secondary" onClick={handleHelp}>
        Need Help?
      </Button>
      <br />
        <br />
      <Typography className={classes.title} color="textPrimary" gutterBottom>
       Enter Websocket server information and telar server gateway
        </Typography>
       
        <TextField
        onChange={handleChange('gateway')}
        label="Gateway"
        value={gateway}
        variant="outlined"
      />
      <br />
      <br />
        <TextField
        InputProps={{
          readonly: true,
        }}
        label="Payload Secret"
        value={payloadSecret}
        variant="outlined"
      />
      <br />
      <br />
        <TextField
        required
        onChange={handleChange('websocketURL')}
        label="Websocket URL"
        value={websocketURL}
        variant="outlined"
      />
      <br />
      <br />
      
        {loadingWebsocket && (<FormControl component="fieldset" className={classes.formControl}>

          <FormGroup>
          <FormControlLabel
          control={checkBox(websocketConnection)}
          checked={websocketConnection}
          label="Check websocket server connection"
        />
           
           
          </FormGroup>
        </FormControl>)}
      </CardContent>
      <CardActions>
      </CardActions>
      <HelpDialog open={helpOpen} onClose={handleCloseHelp} title={'Instruction'}>
       {helpContent()}
      </HelpDialog>
    </Card>
  );
}