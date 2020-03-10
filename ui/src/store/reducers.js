import actionTypes from './actionTypes'

const initialState = {
  inputs: {}
}

  function appReducer(state = initialState, action) {
      const {payload} = action
    switch (action.type) {
      case actionTypes.SET_SETUP_STATE:
        return {...state, setupState: payload.state}
      case actionTypes.SET_SETUP_STEP:
        return {...state, setupStep: payload.step}
      case actionTypes.SET_INPUT:
        return {...state, inputs: {...state.inputs, [payload.key]: payload.value}}
      case actionTypes.SET_DEPLOY_OPEN:
        return {...state, deployOpen: payload.open}
      default:
        return state
    }
  }

export default appReducer