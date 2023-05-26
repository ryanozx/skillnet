import { configureStore } from '@reduxjs/toolkit';
import thunk from 'redux-thunk'; // Optional middleware for handling asynchronous actions
import rootReducer from './reducers/rootReducer';

; // Your root reducer file

const store = configureStore({
  reducer: rootReducer,
  middleware: [thunk],
});

export default store;
