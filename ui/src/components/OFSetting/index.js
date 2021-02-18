import React from 'react';
import { useDispatch, useSelector } from 'react-redux'
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import actions from '../../store/actions'
import TextField from '@material-ui/core/TextField';
import services from '../../services';

const useStyles = makeStyles(theme => ({
    formControl: {
        margin: theme.spacing(3),
    },
    cardActions: {
        justifyContent: "flex-end"
    }
}));

export default function OFSetting() {
  const dispatch = useDispatch()
  const classes = useStyles();
  const projectDirectory = useSelector(state => state['inputs']['projectDirectory'])

    const handleInputChange = (name) => (event) => {
        const {value} = event.currentTarget
        dispatch(actions.setInput(name, value))
      }

    const handleRemoveSocialFromCluster = event => {
        services.dispatchServer(actions.removeSocialFromCluster(projectDirectory))
    };


    return (
        <Card className={classes.root}>
           
            <CardContent>
                <Typography gutterBottom variant="h5" component="h2">
                    Settings
          </Typography>
                
          <br/>

          <Typography className={classes.title} color="textPrimary" gutterBottom>
          Enter your project directory that includes `setup.yml` file
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
            <CardActions className={classes.cardActions}>
                <Button onClick={handleRemoveSocialFromCluster} size="small" color="primary">
                    Remove Social Network from your cluster
      </Button>
            </CardActions>
        </Card>
    );
}