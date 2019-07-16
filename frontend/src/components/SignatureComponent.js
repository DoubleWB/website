import React from 'react';

export default class Signature extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        let options = {
            weekday: 'long',
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
        };
        let styles = {
            backgroundColor: 'purple',
            color: 'white'
        };
        return (
            <tr>
                <td className=".d-sm-none" style={styles}>#{this.props.index + 1}</td>
                <td style={styles}>{this.props.signature.name}</td>
                <td className=".d-sm-none">{new Date(this.props.signature.created_at).toLocaleDateString("en-US", options)}</td>
            </tr>
        );
    }
}
