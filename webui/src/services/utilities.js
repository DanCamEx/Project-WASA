export function isImageSet(url) {
    if (!url || url==""){
        return false;
    }
    return true;
}

export function getDefaultConversationImage(category) {
  if (category == "group") {
    return "/images/default_group_photo.jpg";
  }
  return "/images/default_user_photo.jpg";
}

export function formatTimestamp(timestamp) {
	const date = new Date(timestamp);
	const formattedDate = date.toLocaleString();
	return formattedDate;
}

export function startPolling(method) {
  return setInterval(async () => {
    method();
  }, 10000);
}

export function stopPolling(pollingId) {
  clearInterval(pollingId);
}