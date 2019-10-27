%raw
"import logoSVG from './images/svgs/logo.svg'";

[@react.component]
let make = () => {
  let push = (path: string, _) => ReasonReactRouter.push(path);

  <nav className="container navbar tss-Navbar">
    <div className="navbar-brand">
      <a className="navbar-item" title="The UI of TiDBSlowSQL">
        <img src=[%raw "logoSVG"] alt="TiDBSlowSQL UI" onClick={push("/")} />
      </a>
      <a
        className="navbar-item"
        title="Realtime statistics"
        onClick={push("/")}>
        <i className="fas fa-chart-area fa-lg is-large" />
      </a>
      <a
        className="navbar-item"
        title="Report statistics"
        onClick={push("/report")}>
        <i className="fas fa-file-medical-alt fa-lg is-large" />
      </a>
      <a className="navbar-item">
        <i className="fas fa-cog fa-lg is-large" />
      </a>
    </div>
  </nav>;
};
