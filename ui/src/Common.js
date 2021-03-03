import React from "react";
import { useDispatch, useSelector } from 'react-redux'

import DirectoryDialog from "./components/DirectoryDialog";
import Routes from "./Routes";
import services from "./services";
import actions from './store/actions';

const Common = () => {
  const [directoryOpen, setDirectoryOpen] = React.useState(false);
  const dispatch = useDispatch();

  const handleCloseDir = () => {
    setDirectoryOpen(false);
  };
  const handleOpenDir = () => {
    setDirectoryOpen(true);
  };

  React.useEffect(() => {
    let cookieValue = document.cookie.replace(
      /(?:(?:^|.*;\s*)telar-config-dir\s*\=\s*([^;]*).*$)|^.*$/,
      "$1"
    );
    if (!cookieValue) {
      services.dispatchServer(actions.getProjectDirctory());
    } else {
      dispatch(
        actions.setInput("projectDirectory", cookieValue)
        );
      }
      handleOpenDir();
  }, []);

  return (
  
      <DirectoryDialog open={directoryOpen} onClose={handleCloseDir} />
  
  );
};

export default Common;
