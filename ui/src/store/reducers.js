import actionTypes from './actionTypes'

const initialState = {
  inputs: {},
  infoDialog: {message: "", url: "", open: false}
}

  function appReducer(state = initialState, action) {
      const {payload} = action
    switch (action.type) {
      case actionTypes.SET_SETUP_STATE:
        return {...state, setupState: payload.state}
      case actionTypes.SET_SETUP_STEP:
        return {...state, setupStep: payload.step}
      case actionTypes.SET_STEP_CONDITION:
        return {...state, stepCondition: {[payload.step]: payload.valid}}
      case actionTypes.SET_INPUT:
        return {...state, inputs: {...state.inputs, [payload.key]: payload.value}}
      case actionTypes.SET_DEPLOY_OPEN:
        return {...state, deployOpen: payload.open}
      case actionTypes.SET_SETUP_DEFAULT_VALUES:
        return {...state, inputs: {...state.inputs, ...payload.clientInputs}}
      case actionTypes.POP_MESSAGE:
        return {...state, popMessage: payload.message}
      case actionTypes.SHOW_INFO_DIALOG:
        return {...state, infoDialog: {message: payload.message, url: payload.url, open: true }}
      case actionTypes.HIDE_INFO_DIALOG:
        return {...state, infoDialog: {message: "", url: "", open: false }}
      default:
        return state
    }
  }

export default appReducer