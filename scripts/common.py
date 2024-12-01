


class BColors:
    HEADER = "\033[95m"
    OKBLUE = "\033[94m"
    OKCYAN = "\033[96m"
    OKGREEN = "\033[92m"
    WARNING = "\033[93m"
    FAIL = "\033[91m"
    ENDC = "\033[0m"
    BOLD = "\033[1m"
    UNDERLINE = "\033[4m"


def parse_formatted_number(formatted_num):
    try:
        # Handle strings with 'M' or 'K'
        if isinstance(formatted_num, str):
            if formatted_num.endswith("M"):  # For millions
                return int(
                    float(formatted_num[:-1]) * 1_000_000
                )  # Remove 'M' and convert
            elif formatted_num.endswith("K"):  # For thousands
                return int(float(formatted_num[:-1]) * 1_000)  # Remove 'K' and convert
            elif formatted_num.startswith("€"):  # Remove € if present
                return int(formatted_num[1:])
            else:  # Assume it's a plain number as a string
                return int(float(formatted_num))
        # If already an integer, return as is
        elif isinstance(formatted_num, (int, float)):
            return int(formatted_num)
        else:
            raise ValueError(f"Unexpected format: {formatted_num}")
    except Exception as e:
        print(f"Error parsing number: {formatted_num}, Error: {e}")
        return 0  # Return 0 or handle it as per your requirement


# def convert_to_datetime(time_str: str) -> datetime:
#     """
#     Convert a relative time string like '1d ago', '4h ago', or '4m ago' to a datetime object.
#     """
#     # Get current time 
#     now = datetime.now()

#     # Using regex to analyze the 
#     match = re.match(r"(\d+)([dhm]) ago", time_str)
#     if not match:
#         raise ValueError(f"Invalid time string format: {time_str}")

#     # Analyze the string syntax 
#     value, unit = match.groups()
#     value = int(value)

#     # Conver to time delta based on their unit 
#     if unit == "d":  # days
#         delta = timedelta(days=value)
#     elif unit == "h":  # hours
#         delta = timedelta(hours=value)
#     elif unit == "m":  # minutes
#         delta = timedelta(minutes=value)
#     else:
#         raise ValueError(f"Unsupported time unit: {unit}")

#     # minus the time from the current time 
#     return now - delta
