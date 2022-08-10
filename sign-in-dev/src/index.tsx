import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import "primereact/resources/themes/fluent-light/theme.css";
import "primereact/resources/primereact.min.css";
import "primeicons/primeicons.css";
import 'primeflex/primeflex.css'
import '@mdi/font/css/materialdesignicons.css'

import {
  BrowserRouter,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import SignIn from "./views/sign-in"
import AppPropertiesProvider from './utils/properties';
import ApiProvider from './utils/api-connector';
import { CookiesProvider } from 'react-cookie';
import KeepAlive from './views/keepalive';
import UnAuthorized from './views/unauthorized';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  <CookiesProvider>
    <ApiProvider>
      <BrowserRouter>
        <AppPropertiesProvider>
          <Routes>
            <Route path='*' element={<Navigate to="/sign-in" replace />} />
            <Route path="/sign-in" element={<SignIn />} />
            <Route path="/keepalive" element={<KeepAlive />} />
            <Route path="/unauthorized" element={<UnAuthorized />} />
          </Routes>
        </AppPropertiesProvider>
      </BrowserRouter>
    </ApiProvider>
  </CookiesProvider>

);

