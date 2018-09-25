def update_values_api(values_dict):
    payload = {}
    config = {}
    with open("config.json", "r") as f:
        config = json.load(f)
    headers = { "Authorization" : config["APIAuthorizedToken"] }
    for key in values_dict:
        payload[key] = values_dict[key]
    r = requests.post("http://" + config["APIAddress"] + ":" + str(config["APIPort"]) + "/cpd-update", json=payload, headers=headers)

    for i in range(0,2):
        if r.status_code == requests.codes.ok:
            return True
        else:
            time.sleep(5)
            r = requests.post("http://" + config["APIAddress"] + ":" + str(config["APIPort"]) + "/cpd-update", json=payload, headers=headers)
    return False

import requests, time, json
