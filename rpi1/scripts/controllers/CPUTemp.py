class CPUTemp:
    def __init__(self, temp_file_path="/sys/class/thermal/thermal_zone0/temp"):
        self.temp_file_path = temp_file_path

    def __enter__(self):
        self.open()
        return self

    def open(self):
        """Opens file with rpi temperature from /sys"""
        self.tempfile = open(self.temp_file_path, "r")
    
    def read(self):
        """Reads rpi self temperature from file"""
        self.tempfile.seek(0)
        return self.tempfile.read().rstrip()
    
    def convert_c_to_f(self, c):
        """Converts celsius degrees value to fahrenheit degrees value"""
        return c * 9.0 / 5.0 + 32.0

    def get_temperature_in_c(self):
        """Returns temperature in celsius degrees"""
        temp_raw = self.read()
        return float(temp_raw[:-3] + "." + temp_raw[-3:])

    def get_temperature_in_f(self):
        """Returns temperature in fahrenheit degrees"""
        return self.convert_c_to_f(self.get_temperature_in_c())
    
    def get_temperature(self):
        """Returns temperature (currently in celsius degrees)"""
        return self.get_temperature_in_c()

    def __exit__(self, type, value, traceback):
        self.close()
            
    def close(self):
        """Closes file with rpi temperature from /sys"""
        self.tempfile.close()
