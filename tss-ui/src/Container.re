[@react.component]
let make = () => {
  let url = ReasonReactRouter.useUrl();

  let show =
    switch (url.path) {
    | [] => <Realtime />
    | ["report"] => <Report />
    | _ => <Realtime />
    };

  <div className="container tss-Container"> show </div>;
};
