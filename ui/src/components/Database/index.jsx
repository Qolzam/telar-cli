import React from 'react';
import { makeStyles, withStyles } from '@material-ui/core/styles';
import { useDispatch, useSelector } from 'react-redux'
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

export default function Database() {
  const classes = useStyles();

  const dispatch = useDispatch()
  const mongoDBUsername = useSelector(state => state['inputs']['mongoDBUsername'])
  const mongoDBPassword = useSelector(state => state['inputs']['mongoDBPassword'])
  const mongoDBName = useSelector(state => state['inputs']['mongoDBName'])
  const mongoDBConnection = useSelector(state => state['inputs']['mongoDBConnection'])

  const bull = <span className={classes.bullet}>•</span>;
  const [state, setState] = React.useState({
    gilad: true,
    jason: false,
    antoine: false,
  });
  const [helpOpen, setHelpOpen] = React.useState(false);

  // 0 : instruction
  // 1 : warning
  const [helpType, setHelpType] = React.useState(0);

  const handleHelp = (helpType) => {
    setHelpOpen(true);
    setHelpType(helpType)
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

  const warningContent = () => (
    <Card className={classes.root} variant="outlined">
      <CardContent>
        <Typography variant={'h5'} color="textPrimary" gutterBottom>
        ⚠️ We strictly recommend you to use your own MongoDB. 
        Public MongoDB only use for testing, because everybody can see your data and also there is no garanty for data persistance. 
        Public MongoDB is a free and shared MongoDB account! 
  </Typography>
        <br />
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
        <Button variant="outlined" color="secondary" onClick={() => handleHelp(0)}>
          Need Help?
      </Button>
        {'  '}
        <Button variant="outlined" color="primary" onClick={() => handleHelp(1)}>
          Use Public MongoDB
      </Button>
      <br />
        <br />
      <Typography className={classes.title} color="textPrimary" gutterBottom>
       Fill the field for MongoDB database
        </Typography>
        <br />
        <TextField
        required
        id="outlined-required"
        label="Username"
        
        value={mongoDBUsername}
        onChange={handleChange('mongoDBUsername')}
        variant="outlined"
      />
      <br />
      <br />
        <TextField
        required
        id="outlined-required"
        label="Password"
        
        value={mongoDBPassword}
        onChange={handleChange('mongoDBPassword')}
        variant="outlined"
      />
      <br />
      <br />
        <TextField
        required
        id="outlined-required"
        label="Database Name"
        value={mongoDBName}
        onChange={handleChange('mongoDBName')}
        
        variant="outlined"
      />
      <br />
      <br />
      
        <FormControl component="fieldset" className={classes.formControl}>

          <FormGroup>
          
          <FormControlLabel
              control={checkBox(mongoDBConnection)}
              checked={mongoDBConnection}
              label="Check MongoDB connection"
            />

           
          </FormGroup>
        </FormControl>
      </CardContent>
      <CardActions>
      </CardActions>
      <HelpDialog open={helpOpen} onClose={handleCloseHelp} title={helpType === 0 ? 'Instruction' : '⚠️Warning ⚠️'}>
        {helpType === 0 ? helpContent() : warningContent()}
      </HelpDialog>
    </Card>
  );
}