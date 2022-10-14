import React from 'react';
import ReactDOM from 'react-dom';

import { App } from './App';

const root = document.getElementById('app');

if (root !== null) {
    ReactDOM.render(<App />, root);
    module.hot.accept();
}
