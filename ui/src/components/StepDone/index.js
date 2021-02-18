import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { useDispatch, useSelector } from 'react-redux'
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from '@material-ui/core/CardMedia';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import services from '../../services';

const useStyles = makeStyles({
  root: {
    maxWidth: 645,
  },
  media: {
    height: 277,
  },
});

export default function StepDone() {
  const classes = useStyles();
  const socialDomain = useSelector(state => state['inputs']['socialDomain'])

  const handleLearnMore = () => {
    services.openURL("https://telar.dev")
  }

  return (
    <Card className={classes.root}>
      <CardActionArea>
        <CardMedia
          className={classes.media}
          image="https://miro.medium.com/max/5748/1*M_jDiQwGmGvrKH_H1peKWQ.png"
          title="We are done!"
        />
        <CardContent>
          <Typography gutterBottom variant="h5" component="h2">
            Done
          </Typography>
          <Typography variant="body2" color="textSecondary">
          Wait until the build and deployment is done. To start login into admin page {`${socialDomain}/admin/login`}. Wait until you seen welcome page.
          Now your Telar social is ready to use!😍🏆
          - Signup page: {`${socialDomain}/auth/signup`}
          - Login page: {`${socialDomain}/auth/login`}
          </Typography>
        </CardContent>
      </CardActionArea>
      <CardActions>
        <Button size="small" onClick={handleLearnMore} color="primary">
          Learn More
        </Button>
      </CardActions>
    </Card>
  );
}