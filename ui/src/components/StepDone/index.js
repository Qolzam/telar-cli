import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from '@material-ui/core/CardMedia';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';

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

  return (
    <Card className={classes.root}>
      <CardActionArea>
        <CardMedia
          className={classes.media}
          image="https://miro.medium.com/max/5748/1*XOOHjCzszsqjjB91e-wsyg.png"
          title="We are done!"
        />
        <CardContent>
          <Typography gutterBottom variant="h5" component="h2">
            Done
          </Typography>
          <Typography variant="body2" color="textSecondary" component="p">
          Now your Telar social is ready to use!😍🏆

          </Typography>
        </CardContent>
      </CardActionArea>
      <CardActions>
        <Button size="small" href={'https://telar.press'} target={'_blank'} color="primary">
          Learn More
        </Button>
      </CardActions>
    </Card>
  );
}