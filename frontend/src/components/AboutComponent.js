import React from 'react';

export default class About extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <p className="excerpt">
                This website is project and experiment to practice hosting and maintaining a live service.
                <br/>
                <br/>
                As time goes on, I'd like to push this website forward into more and more clean
                code and deployment practices. As this continues, having a home for hosting information
                about my projects and other going-ons is another benefit of this good practice!
                <br/>
                <br/>
                Currently this site is mainly a landing page for my github and linkedin pages,
                but feel free to leave a signature if you'd like so I can see the variety of people
                that visit!
                <br/>
                <br/>
                Look forward to other improvements along the way.
            </p>
        );
    }
}
