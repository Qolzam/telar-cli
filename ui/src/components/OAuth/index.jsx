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

export default function OAuth() {
  const classes = useStyles();

  const dispatch = useDispatch()
  const githubOAuthSecret = useSelector(state => state['inputs']['githubOAuthSecret'])
  const bull = <span className={classes.bullet}>•</span>;
  const [helpOpen, setHelpOpen] = React.useState(false);

  const handleHelp = () => {
    setHelpOpen(true);
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

  return (
    <Card className={classes.root} variant="outlined">
      <CardContent>
        <Typography className={classes.title} color="textPrimary" gutterBottom>
         Enter your valid Githup OAuth secret key.
        </Typography>
        <Button variant="outlined" color="secondary" onClick={handleHelp}>
        Need Help?
      </Button>
        <br />
      <br />
        <TextField
        required
        id="outlined-required"
        label="Github OAuth Secret key"
        
        variant="outlined"
        value={githubOAuthSecret}
        onChange={handleChange('githubOAuthSecret')}
      />
      <br />
      <br />
      </CardContent>
      <CardActions>
      </CardActions>
      <HelpDialog open={helpOpen} onClose={handleCloseHelp} title={'Instruction'}>
       {helpContent()}
      </HelpDialog>
    </Card>
  );
}