import React from 'react';
import { useSelector } from 'react-redux'
import { makeStyles } from '@material-ui/core/styles';
import { SelectSetup } from '../../components'
import { OFCC } from '../../components'
import Container from '@material-ui/core/Container';
import StepDone from '../../components/StepDone';

const useStyles = makeStyles(theme => ({
  root: {
    '& .MuiTextField-root': {
      margin: theme.spacing(1),
      width: 200,
    },
  },
}));

export default function SetupComponent() {
  const setupState = useSelector(state => state['setupState'])
  const [, setValue] = React.useState('Controlled');


  return (
    <>
      <div style={{marginTop: 20}}></div>
 <Container maxWidth="sm">
     { setupState === 'init' ? <SelectSetup />
     : (setupState === 'progress' ? <OFCC />  :  <StepDone/>)}
      </Container>
    </>
  );
}