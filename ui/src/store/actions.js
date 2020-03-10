import actionTypes from './actionTypes'



const setSetupState = (state) => {
    return {
        type: actionTypes.SET_SETUP_STATE,
        payload: {state}
    }
}

const setSetupStep = (step) => {
    return {
        type: actionTypes.SET_SETUP_STEP,
        payload: {step}
    }
}

const setInput = (key,value) => {
    return {
        type: actionTypes.SET_INPUT,
        payload: {key, value}
    }
}

const setDeployOpen = (open) => {
    return {
        type: actionTypes.SET_DEPLOY_OPEN,
        payload: {open}
    }
}

export default {
    setSetupState,
    setSetupStep,
    setInput,
    setDeployOpen
}