import React from 'react';
import HomeComponent from './components/HomeComponent'
import AboutComponent from './components/AboutComponent'
import { BrowserRouter as Router, Link, Route } from 'react-router-dom'

export default class App extends React.Component {
    render() {
        return (
            <Router>
                <div className="container">
                <h1>Will B. </h1>
                <h2>Robotics // AI // Devops </h2>
                    <Link to="/">Home</Link> 
                    &nbsp;
                    <Link to="/about">About</Link> 
                    &nbsp;
                    <p></p>
                    <Route
                        path="/"
                        render={() => <HomeComponent />} />
                    <Route
                        exact path="/about"
                        render={() => <AboutComponent />} />
                </div>
            </Router>
        )
    }
}
