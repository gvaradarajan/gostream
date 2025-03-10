import { grpc } from "@improbable-eng/grpc-web";
import { AddStreamRequest, AddStreamResponse, ListStreamsRequest, ListStreamsResponse } from "./gen/proto/stream/v1/stream_pb";
import { StreamServiceClient, ServiceError } from "./gen/proto/stream/v1/stream_pb_service";
import { dialWebRTC } from "@viamrobotics/rpc";

const signalingAddress = `${window.location.protocol}//${window.location.host}`;
const host = "local";

async function startup() {
	const webRTCConn = await dialWebRTC(signalingAddress, host);
	const streamClient = new StreamServiceClient(host, { transport: webRTCConn.transportFactory });

	let pResolve: (value: string[]) => void;
	let pReject: (reason?: any) => void;
	let namesPromise = new Promise<string[]>((resolve, reject) => {
		pResolve = resolve;
		pReject = reject;
	});
	const listRequest = new ListStreamsRequest();
	streamClient.listStreams(listRequest, (err: ServiceError, resp: ListStreamsResponse) => {
		if (err) {
			pReject(err);
			return
		}
		pResolve(resp.getNamesList());
	});
	const names = await namesPromise;

	webRTCConn.peerConnection.ontrack = async event => {
		const video = document.createElement('video');
		video.srcObject = event.streams[0];
		video.autoplay = true;
		video.controls = false;
		video.playsInline = true;
		const streamName = event.streams[0].id;
		const streamContainer = document.getElementById(`stream-${streamName}`)!;
		streamContainer.getElementsByTagName("button")[0].remove();
		streamContainer.appendChild(video);
	}

	for (const name of names) {
		const container = document.createElement("div");
		container.id = `stream-${name}`;
		const button = document.createElement("button");
		button.innerText = `Start ${name}`
		button.onclick = async (e) => {
			e.preventDefault();

			button.disabled = true;

			const addRequest = new AddStreamRequest();
			addRequest.setName(name);
			streamClient.addStream(addRequest, (err: ServiceError, resp: AddStreamResponse) => {
				if (err) {
					console.error(err);
					button.disabled = false;
				}
			});
		}
		container.appendChild(button);
		document.body.appendChild(container);
	}
}
startup().catch(e => {
	console.error(e);
});
