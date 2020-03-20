import React from 'react';
import { useDispatch, useSelector } from 'react-redux'
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import actions from '../../store/actions';
import services from '../../services';

export default function DialogInfo() {
  const dispatch = useDispatch()
  const open = useSelector(state => state.infoDialog.open)
  const message = useSelector(state => state.infoDialog.message)
  const url = useSelector(state => state.infoDialog.url)

  const handleClose = () => {
    dispatch(actions.hideInfoDialog())
  };

  const handleOpenURL = (url) => {
    services.openURL(url)
    dispatch(actions.hideInfoDialog())

  };

  return (
    
      <Dialog
        open={open}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">{"Info"}</DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
           {message && `${message[0].toUpperCase()}${message.slice(1)}`}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} color="primary">
            Close
          </Button>
          {(url && url !== "") && <Button onClick={() => handleOpenURL(url)} color="primary" autoFocus>
            Instruction
          </Button>}
        </DialogActions>
      </Dialog>
  );
}