import { h, Fragment } from "preact";
import { useEffect, useState } from "preact/hooks";
import api from "../../api";

interface AlertViewState {
  alerts: Alert[];
  error?: string;
  loading: boolean;
}

interface ErrorViewProps {
  error: string;
}

const ErrorView = (props: ErrorViewProps) => {
  return <div>{props.error}</div>;
};

interface SuccessViewProps {
  alerts: Alert[];
}

const SuccessView = (props: SuccessViewProps) => {
  return <>
    {props.alerts.map((alert) => {
        return (
          <div>
            <div>
              {alert.id} {alert.labels["alertname"]}
            </div>
            <div>
              {Object.keys(alert.labels).map((key) => {
                return (
                  <span>
                    {key}: {alert.labels[key]}
                  </span>
                );
              })}
            </div>
          </div>
        );
      })}
    </>
};

const AlertView = () => {
  const [alerts, setAlerts] = useState<AlertViewState>({
    alerts: [],
    loading: true,
  });

  const fetchAlerts = async () => {
    await api
      .getAlerts()
      .then((newAlerts) => {
        setAlerts({
          ...alerts,
          alerts: newAlerts,
          loading: false,
        });
      })
      .catch((error) => {
        setAlerts({
          ...alerts,
          error: error.toString(),
          loading: false,
        });
      });
  };

  useEffect(() => {
    if (alerts.loading) {
      fetchAlerts();
    }
  }, [alerts]);

  if (alerts.loading) {
    return <div>Loading...</div>;
  } else if (alerts.error) {
    return <ErrorView error={alerts.error} />;
  } else {
    return <SuccessView alerts={alerts.alerts} />;
  }
};

export default AlertView;
