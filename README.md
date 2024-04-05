# pjsftools
Util tools for working with Salesforce.

Requires `sf` (the salesforce command line tool) to work, already authenticated with the org(s) you want to work with.

## pjsftools getperms

This gets the permissions for a field by showing the profile ids and the permssions for that field.

    $ ./pjsftools getperms -org trailheadLwc Contact.Title 
    getperms using org trailheadLwc
    0PS2w000002jK8aGAE:RW;0PS2w000002jK8ZGAU:RW;0PS2w000002jK8EGAU:RW;0PS2w000002jK83GAE:RW;0PS2w00000AZ9I2GAL:RW;0PS2w000002jK8XGAU:RW;0PS2w000002jK8WGAU:RW;0PS2w000002jK8VGAU:RW;0PS2w000002jK8BGAU:RW;0PS2w000002jK8fGAE:RW;0PS2w000002jK86GAE:RW;0PS2w000002jK85GAE:RW;0PS2w000002jK8CGAU:RW;0PS2w000002jK8DGAU:RW;0PS2w000002jK8OGAU:RW;0PS2w000002jK8IGAU:RW;0PS2w000002jK8TGAU:RW;0PS2w000002jK8GGAU:RW;0PS2w000002jK8MGAU:RW;0PS2w000004kHReGAM:RW;0PS2w000002jK8QGAU:RW;0PS2w000002jK8NGAU:RW;0PS2w000002jK8UGAU:RW;0PS2w000002jK8LGAU:RW;0PS2w000002jK8KGAU:RW;0PS2w000002jK82GAE:RW;0PS2w000002jK8HGAU:RW;0PS2w000002jK89GAE:RW;0PS2w000002jK8PGAU:RW;0PS2w000002jK8SGAU:RW;0PS2w000002jK8RGAU:RW;0PS2w000002jK87GAE:RW;0PS2w000008Xdc4GAC:RW;0PS2w000002jK8JGAU:RW;0PS2w000002jK88GAE:RW;0PS2w000002ZO3PGAW:RW;0PS2w000002jK84GAE:RW;0PS2w000002jK8AGAU:RW;0PS2w000002jK81GAE:RW;0PS2w000004kpcXGAQ:RW;0PS2w000002ZO3QGAW:RW;0PS2w000002jK8FGAU:RW

It can use names instead of profile ID

    $ ./pjsftools getperms -org trailheadLwc -names Contact.Title
    getperms using org trailheadLwc
    Analytics+Cloud+Integration+User:RW;Analytics+Cloud+Security+User:RW;Authenticated+Website:RW;Authenticated+Website:RW;B2B+Reordering+Portal+Buyer+Profile:RW;Chatter+External+User:RW;Chatter+Free+User:RW;Chatter+Moderator+User:RW;Contract+Manager:RW;Cross+Org+Data+Proxy+User:RW;Custom%3A+Marketing+Profile:RW;Custom%3A+Sales+Profile:RW;Custom%3A+Support+Profile:RW;Customer+Community+Login+User:RW;Customer+Community+Plus+Login+User:RW;Customer+Community+Plus+User:RW;Customer+Community+User:RW;Customer+Portal+Manager+Custom:RW;Customer+Portal+Manager+Standard:RW;External+Apps+Login+User:RW;External+Identity+User:RW;Force.com+-+App+Subscription+User:RW;Force.com+-+Free+User:RW;Gold+Partner+User:RW;High+Volume+Customer+Portal:RW;High+Volume+Customer+Portal+User:RW;Identity+User:RW;Marketing+User:RW;Partner+App+Subscription+User:RW;Partner+Community+Login+User:RW;Partner+Community+User:RW;Read+Only:RW;Salesforce+API+Only+System+Integrations:RW;Silver+Partner+User:RW;Solution+Manager:RW;Standard+Guest:RW;Standard+Platform+User:RW;Standard+User:RW;System+Administrator:RW;System+Administrator+v2:RW;Testing+1+Profile:RW;Work.com+Only+User:RW

If org is not specified, it uses the default org, as per the `sf` command.

If a field has no permissons for a profile, it may not show up (as there is no record there).

## pjsftools setperms

This is the reverse of getperms, and can set the permissions for a field.

    $ ./pjsftools setperms -org trailheadLwc Contact.Title "0PS2w000002jK8aGAE:RW;0PS2w000002jK8ZGAU:RW"
    setperms using org trailheadLwc
    2 permissions updated

It can also take names.

    $ ./pjsftools setperms -org trailheadLwc -names Contact.Title "Analytics+Cloud+Integration+User:RW;Analytics+Cloud+Security+User:RW"
    setperms using org trailheadLwc
    2 permissions updated

It can also accept piped input.

    $ ./pjsftools getperms -org trailheadLwc Contact.Title | ./pjsftools setperms -org trailheadLwc Contact.Title         
    Reading from pipe
    getperms using org trailheadLwc
    setperms using org trailheadLwc
    42 permissions updated

## Misc

The format for the values are:

* N = Not Visible (read only is not applicable).
* R = Visible, Read Only.
* RW = Visible, Not read only (Writable).
