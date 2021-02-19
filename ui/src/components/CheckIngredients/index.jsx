import React from 'react';
import { makeStyles, withStyles } from '@material-ui/core/styles';
import { useDispatch, useSelector } from 'react-redux'
import { green } from '@material-ui/core/colors';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import FormLabel from '@material-ui/core/FormLabel';
import FormControl from '@material-ui/core/FormControl';
import FormGroup from '@material-ui/core/FormGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';
import CircularProgress from '@material-ui/core/CircularProgress';
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

export default function CheckIngredients() {
  const dispatch = useDispatch()


  const loadingCheckIngredients = useSelector(state => state['inputs']['loadingCheckIngredients'])

  const cloneTelarWeb = useSelector(state => state['inputs']['cloneTelarWeb'])
  const cloneTsServerless = useSelector(state => state['inputs']['cloneTsServerless'])
  const cloneTsUi = useSelector(state => state['inputs']['cloneTsUi'])

  const classes = useStyles();

  const [helpOpen, setHelpOpen] = React.useState(false);

  const handleHelp = () => {
    services.openURL("https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/3.md")
  };

  const handleCloseHelp = () => {
    setHelpOpen(false);
  };

  const helpContent = () => (
    <Card className={classes.root} variant="outlined">
    <CardContent>
      <Typography variant={'h5'} color="textPrimary" gutterBottom>
      This step check the ingredients!
  </Typography>
  <br/>
    </CardContent>
    <CardActions>
    </CardActions>
  </Card>
  )

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
          <FormLabel component="legend">Prepare ingredients</FormLabel>
          <Typography className={classes.title} color="textPrimary" gutterBottom>
          Clean on next to start preparation
        </Typography>
          <FormGroup>

            {
              loadingCheckIngredients && (
                <>
            <FormControlLabel
              control={checkBox(cloneTelarWeb)}
              checked={cloneTelarWeb}
              label={<span>Clone telar-web repository.</span>}
            />
            <FormControlLabel
              control={checkBox(cloneTsServerless)}
              checked={cloneTsServerless}
              label={<span>Clone ts-serverless repository.</span>}
              />
            <FormControlLabel
              control={checkBox(cloneTsUi)}
              checked={cloneTsUi}
              label={<span>Clone ts-ui repository.</span>}
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