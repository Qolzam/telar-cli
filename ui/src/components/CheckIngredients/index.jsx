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

export default function CheckIngredients() {
  const dispatch = useDispatch()

  const githubUsername = useSelector(state => state['inputs']['githubUsername'])

  const loadingCheckIngredients = useSelector(state => state['inputs']['loadingCheckIngredients'])

  const installKubeseal = useSelector(state => state['inputs']['installKubeseal'])
  const githubUsernameRegisterd = useSelector(state => state['inputs']['githubUsernameRegisterd'])
  const cloneTelarWeb = useSelector(state => state['inputs']['cloneTelarWeb'])
  const cloneTsServerless = useSelector(state => state['inputs']['cloneTsServerless'])
  const cloneTsUi = useSelector(state => state['inputs']['cloneTsUi'])
  const openFaaSApp = useSelector(state => state['inputs']['openFaaSApp'])
  const openFaaSAppHasRepos = useSelector(state => state['inputs']['openFaaSAppHasRepos'])
  const githubSSHKey = useSelector(state => state['inputs']['githubSSHKey'])

  const classes = useStyles();
  const bull = <span className={classes.bullet}>•</span>;
  const [state, setState] = React.useState({
    gilad: true,
    jason: false,
    antoine: false,
  });

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
    dispatch(actions.setInput(name,event.currentTarget.checked))
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
      <br/>
      <br/>
        <FormControl component="fieldset" className={classes.formControl}>
          <FormLabel component="legend">Ingredients</FormLabel>
          <FormGroup>
           
          
            <FormControlLabel
              control={
                <Checkbox readOnly={loadingCheckIngredients} checked={openFaaSApp} onChange={handleChange('openFaaSApp')} />
              }
              label={`The OpenFaaS GitHub App is installed.`}
            />
            <FormControlLabel
              control={
                <Checkbox readOnly={loadingCheckIngredients} checked={openFaaSAppHasRepos} onChange={handleChange('openFaaSAppHasRepos')} />
              }
              label={`Add repositories in the OpenFaaS GitHub App.`}
            />
            <FormControlLabel
              control={
                <Checkbox readOnly={loadingCheckIngredients} checked={githubSSHKey} onChange={handleChange('githubSSHKey')} />
              }
              label={`Adding the SSH key to the Github.`}
            />

            {
              loadingCheckIngredients && (
                <>
                <FormControlLabel
              control={checkBox(installKubeseal)}
              checked={installKubeseal}
              label="Install Kubeseal"
            />

            <FormControlLabel
              control={checkBox(githubUsernameRegisterd)}
              checked={githubUsernameRegisterd}
              label={<span>Check <b>{githubUsername}</b> username is registered in OpenFaaS Cloud Community Cluster</span>}
            />
            <FormControlLabel
              control={checkBox(cloneTelarWeb)}
              checked={cloneTelarWeb}
              label={<span>Check the fork on <b>{githubUsername}</b> & Clone telar-web repository.</span>}
            />
            <FormControlLabel
              control={checkBox(cloneTsServerless)}
              checked={cloneTsServerless}
              label={<span>Check the fork on <b>{githubUsername}</b> & Clone ts-serverless repository.</span>}
              />
            <FormControlLabel
              control={checkBox(cloneTsUi)}
              checked={cloneTsUi}
              label={<span>Check the fork on <b>{githubUsername}</b> & Clone ts-ui repository.</span>}
            />
                </>
              )
            }
          </FormGroup>
        </FormControl>
      </CardContent>
      <CardActions>
      </CardActions>
      <HelpDialog open={helpOpen} onClose={handleCloseHelp} title={'Instruction'}>
       {helpContent()}
      </HelpDialog>
    </Card>
  );
}