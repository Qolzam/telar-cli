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

export default function OFCInfo() {
  const classes = useStyles();
  const bull = <span className={classes.bullet}>•</span>;
  const dispatch = useDispatch()
  const appID = useSelector(state => state['inputs']['appID'])
  const ofUsername = useSelector(state => state['inputs']['ofUsername'])
  const ofGateway = useSelector(state => state['inputs']['ofGateway'])
  const socialDomain = useSelector(state => state['inputs']['socialDomain'])
  const secretName = useSelector(state => state['inputs']['secretName'])
  const namespace = useSelector(state => state['inputs']['namespace'])
  const kubeconfigPath = useSelector(state => state['inputs']['kubeconfigPath'])
  const handleInputChange = (name) => (event) => {
    const {value} = event.currentTarget
    dispatch(actions.setInput(name, value))
  }



  const [helpOpen, setHelpOpen] = React.useState(false);

  const handleHelp = () => {
    services.openURL("https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/2.md")
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
          Enter OpenFaaS username. Default is `admin`.
        </Typography>
        <TextField
          required
          id="of-username"
          label="OpenFaaS username"
          onChange={handleInputChange('ofUsername')}
          value={ofUsername}
          variant="outlined"
        />
        <br />
        <br />
        <Typography className={classes.title} color="textPrimary" gutterBottom>
          Enter OpenFaaS gateway e.g. `domain.com` or `localhost.com:31112`
        </Typography>
        <TextField
          required
          id="ofc-domain"
          label="OpenFaaS gateway"
          onChange={handleInputChange('ofGateway')}
          value={ofGateway}
          variant="outlined"
          fullWidth
        />
        <br />
        <br />
        <Typography className={classes.title} color="textPrimary" gutterBottom>
          Enter your social network URL e.g. https://social.com or https://social.com/function
        </Typography>
        <TextField
          required
          id="social-domain"
          label="Social URL"
          onChange={handleInputChange('socialDomain')}
          value={socialDomain}
          variant="outlined"
          fullWidth
        />
        <br />
        <br />
        <Typography className={classes.title} color="textPrimary" gutterBottom>
          Enter OpenFaaS funtion namespace. Default is `openfaas-fn`
        </Typography>
        <TextField
          required
          id="openfaas-namespace"
          label="OpenFaaS namespace"
          onChange={handleInputChange('namespace')}
          value={namespace}
          variant="outlined"
          fullWidth
        />
        <br />
        <br />
        <Typography className={classes.title} color="textPrimary" gutterBottom>
          Enter secret name that will be created on your K8S cluster. Default is `secrets`
        </Typography>
        <TextField
          required
          id="secret-name"
          label="Secret name"
          onChange={handleInputChange('secretName')}
          value={secretName}
          variant="outlined"
          fullWidth
        />
        <br />
        <br />
        <Typography className={classes.title} color="textPrimary" gutterBottom>
          In the case you want to use KUBECONFIG path instead of default context, enter the kube config path. If not let the textbox empty
        </Typography>
        <TextField
          id="kubeconfig-path"
          label="Kubeconfig path"
          onChange={handleInputChange('kubeconfigPath')}
          value={kubeconfigPath}
          variant="outlined"
          fullWidth
        />
      </CardContent>
      <HelpDialog open={helpOpen} onClose={handleCloseHelp} title={'Instruction'}>
       {helpContent()}
      </HelpDialog>
    </Card>
  );
}