import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

import { Home } from './components/Home';
import { Login } from './components/Login';
import { RegisterUser } from './components/Register';

export function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/index.html" element={<Home />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<RegisterUser />} />
      </Routes>
    </BrowserRouter>
  );
}
