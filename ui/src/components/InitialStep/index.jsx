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

export default function InitialStep() {
  const classes = useStyles();
  const bull = <span className={classes.bullet}>•</span>;
  const dispatch = useDispatch()
  const githubUsername = useSelector(state => state['inputs']['githubUsername'])
  const projectDirectory = useSelector(state => state['inputs']['projectDirectory'])
  const handleInputChange = (name) => (event) => {
    const {value} = event.currentTarget
    dispatch(actions.setInput(name, value))
  }



  const [helpOpen, setHelpOpen] = React.useState(false);

  const handleHelp = () => {
    services.openURL("https://github.com/openfaas/community-cluster")
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

  return (
    <Card className={classes.root} variant="outlined">
      <CardContent>
        <Button variant="outlined" color="secondary" onClick={handleHelp}>
        Need Help?
      </Button>
        <br />
        <br />
        <Typography className={classes.title} color="textPrimary" gutterBottom>
          Enter your valid Github user name.
        </Typography>
        <TextField
          required
          id="outlined-required"
          label="Github Username"
          onChange={handleInputChange('githubUsername')}
          value={githubUsername}
          variant="outlined"
        />
        <br />
        <br />
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