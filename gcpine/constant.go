package gcpine

// environment variable key
const (
	EnvKeyChannelSecret      = "CHANNEL_SECRET"
	EnvKeyChannelAccessToken = "CHANNEL_ACCESS_TOKEN"
)

// EventType - Event type Name
type EventType = string

// event Name
const (
	EventTypeTextMessage     EventType = "TextMessage"
	EventTypeImageMessage    EventType = "ImageMessage"
	EventTypeVideoMessage    EventType = "VideoMessage"
	EventTypeAudioMessage    EventType = "AudioMessage"
	EventTypeFileMessage     EventType = "FileMessage"
	EventTypeLocationMessage EventType = "LocationMessage"
	EventTypeStickerMessage  EventType = "StickerMessage"

	EventTypeFollowEvent       EventType = "follow"
	EventTypeUnfollowEvent     EventType = "unfollow"
	EventTypeJoinEvent         EventType = "join"
	EventTypeLeaveEvent        EventType = "leave"
	EventTypeMemberJoinedEvent EventType = "memberJoined"
	EventTypeMemberLeftEvent   EventType = "memberLeft"
	EventTypePostBackEvent     EventType = "postback"
	EventTypeBeaconEvent       EventType = "beacon"
	EventTypeAccountLinkEvent  EventType = "accountLink"
	EventTypeThingsEvent       EventType = "things"
	EventTypeUnsend            EventType = "unsend"
	EventTypeVideoPlayComplete EventType = "videoPlayComplete"
)
