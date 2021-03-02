import React from 'react';
import { useDispatch, useSelector } from 'react-redux'
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import { makeStyles } from '@material-ui/core/styles';
import actions from '../../store/actions'

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
  
export default function DirectoryDialog(props) {
  const classes = useStyles();
  const dispatch = useDispatch()

  const {open, onClose} = props

  const projectDirectory = useSelector(state => state['inputs']['projectDirectory'])
  const handleInputChange = (name) => (event) => {
    const {value} = event.currentTarget
    dispatch(actions.setInput(name, value))
  }

  const handleOK = () => {
    handleClose()
  };

  const handleClose = () => {
    onClose()
  };

  return (
      <Dialog
        open={open}
        disableBackdropClick={true}
        disableEscapeKeyDown={true}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">{"Welcome"}</DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            Please enter your Telar project root directory
          </DialogContentText>
           

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
          
        </DialogContent>
        <DialogActions>
        <Button onClick={handleOK} color="primary">
            OK
          </Button>
        </DialogActions>
      </Dialog>
  );
}