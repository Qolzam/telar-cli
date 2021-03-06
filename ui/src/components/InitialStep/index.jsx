import React from 'react';
import { useDispatch, useSelector } from 'react-redux'
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import actions from '../../store/actions'
import HelpDialog from '../HelpDialog';
import services from '../../services';

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
  address: {
    height: '100%'
  }
});

export default function InitialStep() {
  const classes = useStyles();
  const bull = <span className={classes.bullet}>•</span>;
  const dispatch = useDispatch()
  const appID = useSelector(state => state['inputs']['appID'])
  const appName = useSelector(state => state['inputs']['appName'])
  const companyName = useSelector(state => state['inputs']['companyName'])
  const supportEmail = useSelector(state => state['inputs']['supportEmail'])
  const projectDirectory = useSelector(state => state['inputs']['projectDirectory'])
  const handleInputChange = (name) => (event) => {
    const {value} = event.currentTarget
    dispatch(actions.setInput(name, value))
  }



  const [helpOpen, setHelpOpen] = React.useState(false);

  const handleHelp = () => {
    services.openURL("https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/1.md")
  };

  const handleCloseHelp = () => {
    setHelpOpen(false);
  };
  
  
  const helpContent = () => (
    <Card className={classes.root} variant="outlined">
    <CardContent>
      <Typography variant={'h5'} color="textPrimary" gutterBottom>
      In this step you need to enter your Domain that can access to OpenFaaS gateway and a directory path where you want to keep your project source code. 
  </Typography>
  <br/>
    </CardContent>
    <CardActions>
    </CardActions>
  </Card>
  )

  return (
    <Card className={classes.root} variant="outlined">
      <CardContent>
        <Button variant="outlined" color="secondary" onClick={handleHelp}>
        Need Help?
      </Button>
        <br />
        <br />
        <Typography className={classes.title} color="textPrimary" gutterBottom>
          Enter social network app name. Default is `Telar Social Engine`.
        </Typography>
        <TextField
          required
          id="app-name"
          label="App name"
          onChange={handleInputChange('appName')}
          value={appName}
          variant="outlined"
        />
        <br />
        <br />
        <Typography className={classes.title} color="textPrimary" gutterBottom>
          Enter company name. Default is `Telar`.
        </Typography>
        <TextField
          required
          id="app-name"
          label="Company name"
          onChange={handleInputChange('companyName')}
          value={companyName}
          variant="outlined"
        />
        <br />
        <br />
        <Typography className={classes.title} color="textPrimary" gutterBottom>
          Enter support email. Default is `support@telar.dev`.
        </Typography>
        <TextField
          required
          id="support-email"
          label="Support email"
          onChange={handleInputChange('supportEmail')}
          value={supportEmail}
          variant="outlined"
        />
        <br />
        <br />
        <Typography className={classes.title} color="textPrimary" gutterBottom>
          Enter social network identifier. Default is `my-social-network`.
        </Typography>
        <TextField
          required
          id="social-id"
          label="Scoial Network ID"
          onChange={handleInputChange('appID')}
          value={appID}
          variant="outlined"
        />
        <br />
        <br />
        <Typography className={classes.title} color="textPrimary" gutterBottom>
          Enter your project directory
        </Typography>
        <TextField
          className={classes.address}
          fullWidth
          required
          id="outlined-required"
          label="Project Directory"
          value={projectDirectory}
          onChange={handleInputChange('projectDirectory')}
          variant="outlined"
        />
      </CardContent>
      <HelpDialog open={helpOpen} onClose={handleCloseHelp} title={'Instruction'}>
       {helpContent()}
      </HelpDialog>
    </Card>
  );
}