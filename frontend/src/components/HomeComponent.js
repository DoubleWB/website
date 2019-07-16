import React from 'react';
import SignatureService from '../services/SignatureService'
import Signature from './SignatureComponent'

const signatureService = SignatureService.getInstance();

export default class Home extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      signatures: []
    };

    this.fetchSignatures = this.fetchSignatures.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  fetchSignatures() {
    signatureService.getSignatures().then(res =>
      this.setState({
        signatures: res
      }))
  }

  handleSubmit() {
    signatureService.createSignature({ name: this.state.new_signature }).then(res =>
      this.setState({
        signatures: res
      }))
  }

  handleChange(event) {
    this.setState({
      new_signature: event.target.value
    })
  }

  componentDidMount() {
    this.fetchSignatures();
  }

  render() {
    return (
      <div className="container">
        <br />
        <samp>
          <p>Computer Science and Cognitive Psychology Student. </p>
          <p>Passionate about solving problems elegantly, and becoming the best at that as I can be. </p>
          <p>Graduating January of 2020.</p>
        </samp>
        <hr className="divider"></hr>
        <div className="row">
          <div className="col-3 left">
            <div className="row">
              <a class="btn btn-info col-sm-12" href="https://www.linkedin.com/in/will-becker-059581105/" role="button">Linkedin</a>
            </div>
            &nbsp;
            <div className="row">
              <a class="btn btn-info col-sm-12" href="https://github.com/DoubleWB/" role="button">Personal Github</a>
            </div>
          </div>
          <div className="col-9 right">
            <h3>Guestbook: </h3>
            <form onSubmit={this.handleSubmit}>
              <label>
                Visitor Name:
                &nbsp;
                  <input type="text" value={this.state.new_signature} onChange={this.handleChange} />
              </label>
              &nbsp;
              <input class="button-success" type="submit" value="Sign Guestbook" onClick={this.handleSubmit} />
            </form>
            <table class="table">
              <thead>
                <tr>
                  <th class=".d-sm-none">Number</th>
                  <th>Signature</th>
                  <th class=".d-sm-none">Signed on</th>
                </tr>
              </thead>
              <tbody>
                {this.state.signatures.map((signature, i) =>
                  <Signature signature={signature} index={i} />
                  , this)}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    );
  }
}