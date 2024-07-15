package discord

const (
	ScopeActivitiesRead                        = "activities.read"                          // allows your app to fetch data from a user's "Now Playing/Recently Played" list — not currently available for apps
	ScopeActivitiesWrite                       = "activities.write"                         // allows your app to update a user's activity - not currently available for apps (NOT REQUIRED FOR GAMESDK ACTIVITY MANAGER)
	ScopeApplicationsBuildsRead                = "applications.builds.read"                 // allows your app to read build data for a user's applications
	ScopeApplicationsBuildsUpload              = "applications.builds.upload"               // allows your app to upload/update builds for a user's applications - requires Discord approval
	ScopeApplicationsCommands                  = "applications.commands"                    // allows your app to add commands to a guild - included by default with the bot scope
	ScopeApplicationsCommandsUpdate            = "applications.commands.update"             // allows your app to update its commands using a Bearer token - client credentials grant only
	ScopeApplicationsCommandsPermissionsUpdate = "applications.commands.permissions.update" // allows your app to update permissions for its commands in a guild a user has permissions to
	ScopeApplicationsEntitlements              = "applications.entitlements"                // allows your app to read entitlements for a user's applications
	ScopeApplicationsStoreUpdate               = "applications.store.update"                // allows your app to read and update store data (SKUs, store listings, achievements, etc.) for a user's applications
	ScopeBot                                   = "bot"                                      // for oauth2 bots, this puts the bot in the user's selected guild by default
	ScopeConnections                           = "connections"                              // allows /users/@me/connections to return linked third-party accounts
	ScopeDMChannelsRead                        = "dm_channels.read"                         // allows your app to see information about the user's DMs and group DMs - requires Discord approval
	ScopeEmail                                 = "email"                                    // enables /users/@me to return an email
	ScopeGDMJoin                               = "gdm.join"                                 // allows your app to join users to a group dm
	ScopeGuilds                                = "guilds"                                   // allows /users/@me/guilds to return basic information about all of a user's guilds
	ScopeGuildsJoin                            = "guilds.join"                              // allows /guilds/{guild.id}/members/{user.id} to be used for joining users to a guild
	ScopeGuildsMembersRead                     = "guilds.members.read"                      // allows /users/@me/guilds/{guild.id}/member to return a user's member information in a guild
	ScopeIdentify                              = "identify"                                 // allows /users/@me without email
	ScopeMessagesRead                          = "messages.read"                            // for local rpc server api access, this allows you to read messages from all client channels (otherwise restricted to channels/guilds your app creates)
	ScopeRelationshipsRead                     = "relationships.read"                       // allows your app to know a user's friends and implicit relationships - requires Discord approval
	ScopeRoleConnectionsWrite                  = "role_connections.write"                   // allows your app to update a user's connection and metadata for the app
	ScopeRPC                                   = "rpc"                                      // for local rpc server access, this allows you to control a user's local Discord client - requires Discord approval
	ScopeRPCActivitiesWrite                    = "rpc.activities.write"                     // for local rpc server access, this allows you to update a user's activity - requires Discord approval
	ScopeRPCNotificationsRead                  = "rpc.notifications.read"                   // for local rpc server access, this allows you to receive notifications pushed out to the user - requires Discord approval
	ScopeRPCVoiceRead                          = "rpc.voice.read"                           // for local rpc server access, this allows you to read a user's voice settings and listen for voice events - requires Discord approval
	ScopeRPCVoiceWrite                         = "rpc.voice.write"                          // for local rpc server access, this allows you to update a user's voice settings - requires Discord approval
	ScopeVoice                                 = "voice"                                    // allows your app to connect to voice on user's behalf and see all the voice members - requires Discord approval
	ScopeWebhookIncoming                       = "webhook.incoming"                         // this generates a webhook that is returned in the oauth token response for authorization code grants
)
