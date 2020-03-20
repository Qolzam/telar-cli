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

export default function GoogleReCaptcha() {
  const classes = useStyles();
  const dispatch = useDispatch()
  const siteKeyRecaptcha = useSelector(state => state['inputs']['siteKeyRecaptcha'])
  const recaptchaKey = useSelector(state => state['inputs']['recaptchaKey'])
  const bull = <span className={classes.bullet}>•</span>;
  const [helpOpen, setHelpOpen] = React.useState(false);

  // 0 : instruction
  // 1 : warning
  const [helpType, setHelpType] = React.useState(0);

  const handleHelp = (helpType) => {
    services.openURL("https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/5.md")
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
        <Button variant="outlined" color="secondary" onClick={() => handleHelp(0)}>
          Need Help?
      </Button>
      <br />
        <br />
        <Typography className={classes.title} color="textPrimary" gutterBottom>
        Enter your valid site key and reCaptcha key
        </Typography>
        <TextField
        required
        id="outlined-required"
        label="Site Key"
        
        variant="outlined"
        value={siteKeyRecaptcha}
        onChange={handleChange('siteKeyRecaptcha')}
      />
      <br />
      <br />
        <TextField
        className={classes.address}
        required
        id="outlined-required"
        label="Secret Key"
        
        variant="outlined"
        value={recaptchaKey}
        onChange={handleChange('recaptchaKey')}
      />
      </CardContent>
      <CardActions>
      </CardActions>
      <HelpDialog open={helpOpen} onClose={handleCloseHelp} title={'Instruction'}>
        { helpContent()}
      </HelpDialog>
    </Card>
  );
}