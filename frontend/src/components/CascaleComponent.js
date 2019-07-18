import React from 'react';

export default class Cascale extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <div>
                <div className="centered">
                    <iframe width="560" height="315" src="https://www.youtube.com/embed/DbKkiVTpiSI" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
                </div>
                <p className="excerpt">
                    CaScale is a hardware project that hoped to fill a gap in the high end coffee scale market - a fully smart scale.
                    In the promo video, you can see some intended features, like rotary menus, onboard memory, and recipe execution, which
                    includes segmented recipes, timing, and guidance to improve the consistency of your rate of pouring.
                <br />
                    <br />
                    Other coffee scales allow for some smart feedback to mobile applications through Bluetooth, but our team wanted to build
                    a scale which would actually direct user action, instead of simply providing information and letting the user act independently.
                    In this way, we would hope to minimize the user facing information while maximizing the intelligence of the scale. Further down
                    the pipeline, we envision having a website application to create, network, and share recipes that could be loaded onto the scale.
                    We also wanted to make a scale that maximized precision (up to hundredths of a gram) while minimizing fluctuation, which is a difficulty
                    not many scales have endeavored to tackle.
                <br />
                    <br />
                    CaScale is currently on hiatus due to the founding team members being unable to commit sufficient time to the project,
                    but I hope to one day return to it!
            </p>
            </div>
        );
    }
}
