import actionTypes from './actionTypes'
import Cookies from 'js-cookie'

// ******* Client Actions ******* //
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

const setStepCondition = (step, valid) => {
    return {
        type: actionTypes.SET_STEP_CONDITION,
        payload: {step, valid}
    }
}

const setInput = (key,value) => {
    if (key === 'projectDirectory') {
        Cookies.set('telar-config-dir', value, { expires: 222220 });
    }
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

const setSetupDefaultValues = (payload) => {
    return {
        type: actionTypes.SET_SETUP_DEFAULT_VALUES,
        payload
    }
}

const popMessage = (message) => {
    return {
        type: actionTypes.POP_MESSAGE,
        payload: {message}
    }
}

const showInfoDialog = (message, url) => {
    return {
        type: actionTypes.SHOW_INFO_DIALOG,
        payload: {message, url}
    }
}

const hideInfoDialog = () => {
    return {
        type: actionTypes.HIDE_INFO_DIALOG,
        payload: {}
    }
}


// ******* Server HTTP Actions ******* //
const startStep = () => {
    return {
        type: actionTypes.START_STEP,
        payload: {}
    }
}

const removeSocialFromCluster = (projectDirectory) => {
    return {
        type: actionTypes.REMOVE_SOCIAL_FROM_CLUSTER,
        payload: {projectDirectory}
    }
}

const getProjectDirctory = () => {
    return {
        type: actionTypes.ECHO_PROJECT_DIR,
        payload: {}
    }
}

const checkStep = (payload) => {
    return {
        type: actionTypes.CHECK_STEP,
        payload
    }
}

export default {
    setSetupState,
    setSetupStep,
    setStepCondition,
    setInput,
    setDeployOpen,
    setSetupDefaultValues,
    getProjectDirctory,
    popMessage,
    showInfoDialog,
    hideInfoDialog,
    startStep,
    removeSocialFromCluster,
    checkStep
}