import { LOGIN_REQUEST, LOGIN_SUCCESS, LOGIN_FAILURE } from '../actions/userActions';

const initialState = {
  loading: false,
  isLoggedIn: false,
  user: {},
  error: ''
}

const userReducer = (state = initialState, action) => {
  switch(action.type) {
    case LOGIN_REQUEST:
      return {
        ...state,
        loading: true
      }
    case LOGIN_SUCCESS:
      return {
        loading: false,
        isLoggedIn: true,
        user: action.payload,
        error: ''
      }
    case LOGIN_FAILURE:
      return {
        loading: false,
        isLoggedIn: false,
        user: {},
        error: action.payload
      }
    default:
      return state;
  }
}

export default userReducer;
