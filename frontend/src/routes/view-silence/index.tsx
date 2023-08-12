import { h } from "preact"
import { DefaultService, Silence } from "../../api";
import { useState } from "preact/hooks";
import Loader from "../../components/loader";

const ViewSilence = () => {
    const params = new URLSearchParams(window.location.search);
    const id = params.get("id");

    const [silence, setSilence] = useState<Silence | null>(null);

    const fetchSilence = () => {
        // DefaultService.getSilences()
    };


    return <div>
        <Loader loader={} />
    </div>
}