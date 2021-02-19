import React from 'react';
import { useDispatch, useSelector } from 'react-redux'
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import { Divider } from '@material-ui/core';
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

export default function UserManagement() {
  const classes = useStyles();
  const dispatch = useDispatch()
  const adminUsername = useSelector(state => state['inputs']['adminUsername'])
  const adminPassword = useSelector(state => state['inputs']['adminPassword'])
  const gmail = useSelector(state => state['inputs']['gmail'])
  const gmailPassword = useSelector(state => state['inputs']['gmailPassword'])
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
  const bull = <span className={classes.bullet}>•</span>;

  const handleChange = name => event => {
    dispatch(actions.setInput(name,event.currentTarget.value))
  };


  return (
    <Card className={classes.root} variant="outlined">
      <CardContent>
        <Button variant="outlined" color="secondary" onClick={handleHelp}>
        Need Help?
      </Button>

        <br />
      <br />
        <Typography className={classes.title} color="textPrimary" gutterBottom>
         Enter Admin user information.
        </Typography>
        <TextField
        required
        id="outlined-required"
        label="Admin Email"
        
        variant="outlined"
        value={adminUsername}
        onChange={handleChange('adminUsername')}
      />
      <br />
      <br />
        <TextField
        required
        id="outlined-required"
        label="Password"
        
        variant="outlined"
        value={adminPassword}
        onChange={handleChange('adminPassword')}
      />
      </CardContent>
      <CardContent>
        <Typography className={classes.title} color="textPrimary" gutterBottom>
         Enter your Gmail information.
        </Typography>
        <Typography variant="body2" component="p">
         It will be used for sending confirmation and verification email to users.
          <br />
        </Typography>
        <br />
      <br />
        <TextField
        required
        id="outlined-required"
        label="Email"
        
        variant="outlined"
        value={gmail}
        onChange={handleChange('gmail')}
      />
      <br />
      <br />
        <TextField
        required
        id="outlined-required"
        label="Password"
        variant="outlined"
        value={gmailPassword}
        onChange={handleChange('gmailPassword')}
      />
      </CardContent>
      <CardActions>
      </CardActions>
      <HelpDialog open={helpOpen} onClose={handleCloseHelp} title={'Instruction'}>
       {helpContent()}
      </HelpDialog>
    </Card>
  );
}