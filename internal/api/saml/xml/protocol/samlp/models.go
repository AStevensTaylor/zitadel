// Code generated by https://github.com/gocomply/xsd2go; DO NOT EDIT.
// Models for urn:oasis:names:tc:SAML:2.0:protocol
package samlp

import (
	"encoding/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/saml"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/xenc"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/xml_dsig"
)

// Element
type Extensions struct {
	XMLName xml.Name `xml:"Extensions"`
}

// Element
type Status struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:protocol Status"`

	StatusCode StatusCodeType `xml:"StatusCode"`

	StatusMessage string `xml:"StatusMessage"`

	StatusDetail *StatusDetailType `xml:"StatusDetail"`
}

// Element
type StatusCode struct {
	XMLName xml.Name `xml:"StatusCode"`

	Value string `xml:"Value,attr"`

	StatusCode *StatusCodeType `xml:",any"`
}

// Element
type StatusMessage struct {
	XMLName xml.Name `xml:"StatusMessage"`

	Text string `xml:",chardata"`
}

// Element
type StatusDetail struct {
	XMLName xml.Name `xml:"StatusDetail"`
}

// Element
type AssertionIDRequest struct {
	XMLName xml.Name `xml:"AssertionIDRequest"`

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	AssertionIDRef []string `xml:"AssertionIDRef"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`
}

// Element
type SubjectQuery struct {
	XMLName xml.Name `xml:"SubjectQuery"`

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Subject saml.SubjectType `xml:"Subject"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`
}

// Element
type AuthnQuery struct {
	XMLName xml.Name `xml:"AuthnQuery"`

	SessionIndex string `xml:"SessionIndex,attr,omitempty"`

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	RequestedAuthnContext *RequestedAuthnContextType `xml:"RequestedAuthnContext"`

	Subject saml.SubjectType `xml:"Subject"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`
}

// Element
type RequestedAuthnContext struct {
	XMLName xml.Name `xml:"RequestedAuthnContext"`

	Comparison AuthnContextComparisonType `xml:"Comparison,attr,omitempty"`

	AuthnContextClassRef []string `xml:"AuthnContextClassRef"`

	AuthnContextDeclRef []string `xml:"AuthnContextDeclRef"`
}

// Element
type AttributeQuery struct {
	XMLName xml.Name `xml:"AttributeQuery"`

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Attribute []saml.AttributeType `xml:"Attribute"`

	Subject saml.SubjectType `xml:"Subject"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`
}

// Element
type AuthzDecisionQuery struct {
	XMLName xml.Name `xml:"AuthzDecisionQuery"`

	Resource string `xml:"Resource,attr"`

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Action []saml.ActionType `xml:"Action"`

	Evidence *saml.EvidenceType `xml:"Evidence"`

	Subject saml.SubjectType `xml:"Subject"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`
}

// Element
type AuthnRequest struct {
	XMLName xml.Name `xml:"AuthnRequest"`

	ForceAuthn string `xml:"ForceAuthn,attr,omitempty"`

	IsPassive string `xml:"IsPassive,attr,omitempty"`

	ProtocolBinding string `xml:"ProtocolBinding,attr,omitempty"`

	AssertionConsumerServiceIndex string `xml:"AssertionConsumerServiceIndex,attr,omitempty"`

	AssertionConsumerServiceURL string `xml:"AssertionConsumerServiceURL,attr,omitempty"`

	AttributeConsumingServiceIndex string `xml:"AttributeConsumingServiceIndex,attr,omitempty"`

	ProviderName string `xml:"ProviderName,attr,omitempty"`

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Subject *saml.SubjectType `xml:"Subject"`

	NameIDPolicy *NameIDPolicyType `xml:"NameIDPolicy"`

	Conditions *saml.ConditionsType `xml:"Conditions"`

	RequestedAuthnContext *RequestedAuthnContextType `xml:"RequestedAuthnContext"`

	Scoping *ScopingType `xml:"Scoping"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`
}

// Element
type NameIDPolicy struct {
	XMLName xml.Name `xml:"NameIDPolicy"`

	Format string `xml:"Format,attr,omitempty"`

	SPNameQualifier string `xml:"SPNameQualifier,attr,omitempty"`

	AllowCreate bool `xml:"AllowCreate,attr,omitempty"`
}

// Element
type Scoping struct {
	XMLName xml.Name `xml:"Scoping"`

	ProxyCount int `xml:"ProxyCount,attr,omitempty"`

	IDPList *IDPListType `xml:"IDPList"`

	RequesterID []string `xml:"RequesterID"`
}

// Element
type RequesterID struct {
	XMLName xml.Name `xml:"RequesterID"`

	Text string `xml:",chardata"`
}

// Element
type IDPList struct {
	XMLName xml.Name `xml:"IDPList"`

	IDPEntry []IDPEntryType `xml:"IDPEntry"`

	GetComplete string `xml:"GetComplete"`
}

// Element
type IDPEntry struct {
	XMLName xml.Name `xml:"IDPEntry"`

	ProviderID string `xml:"ProviderID,attr"`

	Name string `xml:"Name,attr,omitempty"`

	Loc string `xml:"Loc,attr,omitempty"`
}

// Element
type GetComplete struct {
	XMLName xml.Name `xml:"GetComplete"`

	Text string `xml:",chardata"`
}

// Element
type Response struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:protocol Response"`

	Id string `xml:"ID,attr"`

	InResponseTo string `xml:"InResponseTo,attr,omitempty"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Issuer *saml.NameIDType `xml:"urn:oasis:names:tc:SAML:2.0:assertion Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Assertion saml.Assertion `xml:"Assertion"`

	Extensions *ExtensionsType `xml:"Extensions"`

	Status StatusType `xml:"Status"`
}

// Element
type ArtifactResolve struct {
	XMLName xml.Name `xml:"ArtifactResolve"`

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Artifact string `xml:"Artifact"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`
}

// Element
type Artifact struct {
	XMLName xml.Name `xml:"Artifact"`

	Text string `xml:",chardata"`
}

// Element
type ArtifactResponse struct {
	XMLName xml.Name `xml:"ArtifactResponse"`

	Id string `xml:"ID,attr"`

	InResponseTo string `xml:"InResponseTo,attr,omitempty"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	Status StatusType `xml:"Status"`
}

// Element
type ManageNameIDRequest struct {
	XMLName xml.Name `xml:"ManageNameIDRequest"`

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	NameID *saml.NameIDType `xml:"NameID"`

	EncryptedID *saml.EncryptedElementType `xml:"EncryptedID"`

	NewID string `xml:"NewID"`

	NewEncryptedID *saml.EncryptedElementType `xml:"NewEncryptedID"`

	Terminate *TerminateType `xml:"Terminate"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`
}

// Element
type NewID struct {
	XMLName xml.Name `xml:"NewID"`

	Text string `xml:",chardata"`
}

// Element
type NewEncryptedID struct {
	XMLName xml.Name `xml:"NewEncryptedID"`

	EncryptedData xenc.EncryptedDataType `xml:"EncryptedData"`

	EncryptedKey []xenc.EncryptedKeyType `xml:"EncryptedKey"`
}

// Element
type Terminate struct {
	XMLName xml.Name `xml:"Terminate"`
}

// Element
type ManageNameIDResponse struct {
	XMLName xml.Name `xml:"ManageNameIDResponse"`

	Id string `xml:"ID,attr"`

	InResponseTo string `xml:"InResponseTo,attr,omitempty"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	Status StatusType `xml:"Status"`
}

// Element
type LogoutRequest struct {
	XMLName xml.Name `xml:"LogoutRequest"`

	Reason string `xml:"Reason,attr,omitempty"`

	NotOnOrAfter string `xml:"NotOnOrAfter,attr,omitempty"`

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	SessionIndex []string `xml:"SessionIndex"`

	BaseID *saml.BaseIDAbstractType `xml:"BaseID"`

	NameID *saml.NameIDType `xml:"NameID"`

	EncryptedID *saml.EncryptedElementType `xml:"EncryptedID"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`
}

// Element
type SessionIndex struct {
	XMLName xml.Name `xml:"SessionIndex"`

	Text string `xml:",chardata"`
}

// Element
type LogoutResponse struct {
	XMLName xml.Name `xml:"LogoutResponse"`

	Id string `xml:"ID,attr"`

	InResponseTo string `xml:"InResponseTo,attr,omitempty"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	Status StatusType `xml:"Status"`
}

// Element
type NameIDMappingRequest struct {
	XMLName xml.Name `xml:"NameIDMappingRequest"`

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	NameIDPolicy NameIDPolicyType `xml:"NameIDPolicy"`

	BaseID *saml.BaseIDAbstractType `xml:"BaseID"`

	NameID *saml.NameIDType `xml:"NameID"`

	EncryptedID *saml.EncryptedElementType `xml:"EncryptedID"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`
}

// Element
type NameIDMappingResponse struct {
	XMLName xml.Name `xml:"NameIDMappingResponse"`

	Id string `xml:"ID,attr"`

	InResponseTo string `xml:"InResponseTo,attr,omitempty"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	Status StatusType `xml:"Status"`
}

// XSD ComplexType declarations

type RequestAbstractType struct {
	XMLName xml.Name

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	InnerXml string `xml:",innerxml"`
}

type ExtensionsType struct {
	XMLName xml.Name

	InnerXml string `xml:",innerxml"`
}

type StatusResponseType struct {
	XMLName xml.Name

	Id string `xml:"ID,attr"`

	InResponseTo string `xml:"InResponseTo,attr,omitempty"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	Status StatusType `xml:"Status"`

	InnerXml string `xml:",innerxml"`
}

type StatusType struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:protocol Status"`

	StatusCode StatusCodeType `xml:"StatusCode"`

	StatusMessage string `xml:"StatusMessage"`

	StatusDetail *StatusDetailType `xml:"StatusDetail"`

	InnerXml string `xml:",innerxml"`
}

type StatusCodeType struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:protocol StatusCode"`

	Value string `xml:"Value,attr"`

	StatusCode *StatusCodeType `xml:",any"`

	InnerXml string `xml:",innerxml"`
}

type StatusDetailType struct {
	XMLName xml.Name

	InnerXml string `xml:",innerxml"`
}

type AssertionIDRequestType struct {
	XMLName xml.Name

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	AssertionIDRef []string `xml:"AssertionIDRef"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	InnerXml string `xml:",innerxml"`
}

type SubjectQueryAbstractType struct {
	XMLName xml.Name

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Subject saml.SubjectType `xml:"Subject"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	InnerXml string `xml:",innerxml"`
}

type AuthnQueryType struct {
	XMLName xml.Name

	SessionIndex string `xml:"SessionIndex,attr,omitempty"`

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	RequestedAuthnContext *RequestedAuthnContextType `xml:"RequestedAuthnContext"`

	Subject saml.SubjectType `xml:"Subject"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	InnerXml string `xml:",innerxml"`
}

type RequestedAuthnContextType struct {
	XMLName xml.Name

	Comparison AuthnContextComparisonType `xml:"Comparison,attr,omitempty"`

	AuthnContextClassRef []string `xml:"AuthnContextClassRef"`

	AuthnContextDeclRef []string `xml:"AuthnContextDeclRef"`

	InnerXml string `xml:",innerxml"`
}

type AttributeQueryType struct {
	XMLName xml.Name

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Attribute []saml.AttributeType `xml:"Attribute"`

	Subject saml.SubjectType `xml:"Subject"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	InnerXml string `xml:",innerxml"`
}

type AuthzDecisionQueryType struct {
	XMLName xml.Name

	Resource string `xml:"Resource,attr"`

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Action []saml.ActionType `xml:"Action"`

	Evidence *saml.EvidenceType `xml:"Evidence"`

	Subject saml.SubjectType `xml:"Subject"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	InnerXml string `xml:",innerxml"`
}

type AuthnRequestType struct {
	XMLName xml.Name

	ForceAuthn string `xml:"ForceAuthn,attr,omitempty"`

	IsPassive string `xml:"IsPassive,attr,omitempty"`

	ProtocolBinding string `xml:"ProtocolBinding,attr,omitempty"`

	AssertionConsumerServiceIndex string `xml:"AssertionConsumerServiceIndex,attr,omitempty"`

	AssertionConsumerServiceURL string `xml:"AssertionConsumerServiceURL,attr,omitempty"`

	AttributeConsumingServiceIndex string `xml:"AttributeConsumingServiceIndex,attr,omitempty"`

	ProviderName string `xml:"ProviderName,attr,omitempty"`

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Subject *saml.SubjectType `xml:"Subject"`

	NameIDPolicy *NameIDPolicyType `xml:"NameIDPolicy"`

	Conditions *saml.ConditionsType `xml:"Conditions"`

	RequestedAuthnContext *RequestedAuthnContextType `xml:"RequestedAuthnContext"`

	Scoping *ScopingType `xml:"Scoping"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	InnerXml string `xml:",innerxml"`
}

type NameIDPolicyType struct {
	XMLName xml.Name

	Format string `xml:"Format,attr,omitempty"`

	SPNameQualifier string `xml:"SPNameQualifier,attr,omitempty"`

	AllowCreate bool `xml:"AllowCreate,attr,omitempty"`

	InnerXml string `xml:",innerxml"`
}

type ScopingType struct {
	XMLName xml.Name

	ProxyCount int `xml:"ProxyCount,attr,omitempty"`

	IDPList *IDPListType `xml:"IDPList"`

	RequesterID []string `xml:"RequesterID"`

	InnerXml string `xml:",innerxml"`
}

type IDPListType struct {
	XMLName xml.Name

	IDPEntry []IDPEntryType `xml:"IDPEntry"`

	GetComplete string `xml:"GetComplete"`

	InnerXml string `xml:",innerxml"`
}

type IDPEntryType struct {
	XMLName xml.Name

	ProviderID string `xml:"ProviderID,attr"`

	Name string `xml:"Name,attr,omitempty"`

	Loc string `xml:"Loc,attr,omitempty"`

	InnerXml string `xml:",innerxml"`
}

type ResponseType struct {
	XMLName xml.Name `xml:"Response"`

	Id string `xml:"ID,attr"`

	InResponseTo string `xml:"InResponseTo,attr,omitempty"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	Status StatusType `xml:"Status"`

	InnerXml string `xml:",innerxml"`
}

type ArtifactResolveType struct {
	XMLName xml.Name

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Artifact string `xml:"Artifact"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	InnerXml string `xml:",innerxml"`
}

type ArtifactResponseType struct {
	XMLName xml.Name

	Id string `xml:"ID,attr"`

	InResponseTo string `xml:"InResponseTo,attr,omitempty"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	Status StatusType `xml:"Status"`

	InnerXml string `xml:",innerxml"`
}

type ManageNameIDRequestType struct {
	XMLName xml.Name

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	NameID *saml.NameIDType `xml:"NameID"`

	EncryptedID *saml.EncryptedElementType `xml:"EncryptedID"`

	NewID string `xml:"NewID"`

	NewEncryptedID *saml.EncryptedElementType `xml:"NewEncryptedID"`

	Terminate *TerminateType `xml:"Terminate"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	InnerXml string `xml:",innerxml"`
}

type TerminateType struct {
	XMLName xml.Name

	InnerXml string `xml:",innerxml"`
}

type LogoutRequestType struct {
	XMLName xml.Name

	Reason string `xml:"Reason,attr,omitempty"`

	NotOnOrAfter string `xml:"NotOnOrAfter,attr,omitempty"`

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	SessionIndex []string `xml:"SessionIndex"`

	BaseID *saml.BaseIDAbstractType `xml:"BaseID"`

	NameID *saml.NameIDType `xml:"NameID"`

	EncryptedID *saml.EncryptedElementType `xml:"EncryptedID"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	InnerXml string `xml:",innerxml"`
}

type NameIDMappingRequestType struct {
	XMLName xml.Name

	Id string `xml:"ID,attr"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	NameIDPolicy NameIDPolicyType `xml:"NameIDPolicy"`

	BaseID *saml.BaseIDAbstractType `xml:"BaseID"`

	NameID *saml.NameIDType `xml:"NameID"`

	EncryptedID *saml.EncryptedElementType `xml:"EncryptedID"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	InnerXml string `xml:",innerxml"`
}

type NameIDMappingResponseType struct {
	XMLName xml.Name

	Id string `xml:"ID,attr"`

	InResponseTo string `xml:"InResponseTo,attr,omitempty"`

	Version string `xml:"Version,attr"`

	IssueInstant string `xml:"IssueInstant,attr"`

	Destination string `xml:"Destination,attr,omitempty"`

	Consent string `xml:"Consent,attr,omitempty"`

	Issuer *saml.NameIDType `xml:"Issuer"`

	Signature *xml_dsig.SignatureType `xml:"Signature"`

	Extensions *ExtensionsType `xml:"Extensions"`

	Status StatusType `xml:"Status"`

	InnerXml string `xml:",innerxml"`
}

// XSD SimpleType declarations

type AuthnContextComparisonType string

const AuthnContextComparisonTypeExact AuthnContextComparisonType = "exact"

const AuthnContextComparisonTypeMinimum AuthnContextComparisonType = "minimum"

const AuthnContextComparisonTypeMaximum AuthnContextComparisonType = "maximum"

const AuthnContextComparisonTypeBetter AuthnContextComparisonType = "better"
