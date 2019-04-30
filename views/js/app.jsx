class App extends React.Component {

  render() {
    return <Signatures />;
  }

}

class Signatures extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      signatures: []
    };

    this.serverRequest = this.serverRequest.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  serverRequest() {
    $.get("www.doublewb.xyz/api/signs", res => {
      this.setState({
        signatures: res
      });
    });
  }

  handleSubmit(event) {
    $.post("www.doublewb.xyz/api/sign", JSON.stringify({name: this.state.new_signature}), res => {
      this.setState({
        signatures: res
      });
    });
  }

  handleChange(event) {
    this.setState({
      new_signature: event.target.value
    })
  }

  componentDidMount() {
    this.serverRequest();
  }

  render() {
    let combinedStyle = {
      display: 'flex',
      flexDirection: 'row'
    };
    let spacerStyle = {
      margin: '30px',
    };
    let dividerStyle = {
      backgroundColor: '#fff',
      borderTop: '2px dashed #8c8b8b'
    }
    return (
      <div className="container">
        <br />
        <h1>Will B. </h1>
        <h2>Robotics // AI // Devops </h2>
        <samp>
          <p>Computer Science and Cognitive Psychology Student. </p>
          <p>Passionate about solving problems elegantly, and becoming the best at that as I can be. </p>
          <p>Graduating January of 2020.</p>
        </samp>
        <hr style={dividerStyle}></hr>
        <div style={spacerStyle}/>
        <div style={combinedStyle}>
          <div>
          <div style={spacerStyle}/>
          <a class="btn btn-info" href="https://www.linkedin.com/in/will-becker-059581105/" role="button">Linkedin</a>
          <div style={spacerStyle}/>
          <a class="btn btn-info" href="https://github.com/DoubleWB/" role="button">Personal Github</a>
          </div>
          <div style={spacerStyle}/>
          <div>
            <h3>Guestbook: </h3>
            <form onSubmit={this.handleSubmit}>
              <label>
                Visitor Name:
                <input type="text" value={this.state.new_signature} onChange={this.handleChange} />
              </label>
              <input type="submit" value="Sign Guestbook" onClick={this.handleSubmit}/>
            </form>
            <div className="container">
              {this.state.signatures.map(function(signature, i) {
                return (
                  <div className="row" >
                    <Signature signature={signature} index={i}/>
                  </div>);
                }, this)}
            </div>
          </div>
        </div>
      </div>
    );
  }
}

class Signature extends React.Component {
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
      <div className="col-xs-4">
        <div className="panel panel-default">
          <div className="panel-heading" style={styles}>
            #{this.props.index + 1}{" "}
            <span className="center">{this.props.signature.name}</span>
          </div>
          <div className="panel-body joke-hld">{new Date(this.props.signature.created_at).toLocaleDateString("en-US", options)}</div>
        </div>
      </div>
    );
  }
}


ReactDOM.render(<App />, document.getElementById("app"));
